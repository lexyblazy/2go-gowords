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
