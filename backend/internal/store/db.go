package store

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type SqlDb struct {
	db *sql.DB
}

func NewSqlDB(dsn string) (*SqlDb, error) {
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	var mode string
	err = db.QueryRow(`PRAGMA journal_mode=WAL;`).Scan(&mode)
	if err != nil {
		return nil, err
	}

	if mode != "wal" {
		return nil, fmt.Errorf("failed to enable WAL mode, got %s", mode)
	}

	_, err = db.Exec(`
		PRAGMA synchronous=NORMAL;
		PRAGMA busy_timeout=5000;
	`)

	if err != nil {
		return nil, err
	}

	return &SqlDb{
		db: db,
	}, nil
}

func (s *SqlDb) getUserByColumn(colName string, value string) (UserEntity, error) {
	var user UserEntity

	query := fmt.Sprintf("select id, username, moniker, password, created_at, recovery_hash from users where %s = ?", colName)

	err := s.db.QueryRow(query, value).Scan(
		&user.ID,
		&user.Username,
		&user.Moniker,
		&user.Password,
		&user.CreatedAt,
		&user.RecoveryHash)

	if err != nil {
		return user, err
	}

	return user, nil
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func (s *SqlDb) GetUserById(id string) (UserEntity, error) {
	return s.getUserByColumn("id", id)

}

func (s *SqlDb) GetUserByUsername(username string) (UserEntity, error) {
	return s.getUserByColumn("username", normalizeUsername(username))
}

func (s *SqlDb) CreateUser(id string, username string, passwordHash string, moniker string, recoveryHash string) (UserEntity, error) {

	var user UserEntity

	err := s.db.QueryRow(`insert into users(id, username, password, moniker, recovery_hash) 
	values (?, ?, ?, ?, ?) returning id, username, moniker `, id, normalizeUsername(username), passwordHash, moniker, recoveryHash).Scan(
		&user.ID,
		&user.Username,
		&user.Moniker,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *SqlDb) UpdatePassword(userId string, passwordHash string) error {
	_, err := s.db.Exec(`update users set password = ? where id = ?`, passwordHash, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqlDb) UpdateUserStats(updates []UserStatsUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
	INSERT INTO user_stats (
	  user_id,
	  wins_count,
	  best_score,
	  games_played,
	  total_score,
	  created_at,
	  updated_at
	)
	SELECT
	  u.id,
	  ?2,
	  ?3,
	  1,
	  ?3,
	  CURRENT_TIMESTAMP,
	  CURRENT_TIMESTAMP
	FROM users u
	WHERE u.id = ?1
	ON CONFLICT(user_id) DO UPDATE SET
	  games_played = user_stats.games_played + 1,
	  wins_count = user_stats.wins_count + excluded.wins_count,
	  best_score = MAX(user_stats.best_score, excluded.best_score),
	  total_score = user_stats.total_score + excluded.total_score,
	  updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, u := range updates {
		winIncrement := 0
		if u.IsWinner {
			winIncrement = 1
		}
		if _, err := stmt.Exec(
			u.UserID,
			winIncrement,
			u.Score,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SqlDb) Close() {
	s.db.Close()
}
