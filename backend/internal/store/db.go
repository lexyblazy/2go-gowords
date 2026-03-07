package store

import (
	"database/sql"
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

func (s *SqlDb) GetUserByUsername(username string) (UserEntity, error) {

	var user UserEntity
	err := s.db.QueryRow("select id, username, moniker, password, created_at  from users where username = ?", username).Scan(&user.ID,
		&user.Username,
		&user.Moniker,
		&user.Password,
		&user.CreatedAt)

	if err != nil {
		return user, err
	}

	return user, nil

}

func (s *SqlDb) CreateUser(id string, username string, password string, moniker string) (UserEntity, error) {

	var user UserEntity

	err := s.db.QueryRow("insert into users(id, username, password, moniker) values (?, ?, ?, ?) returning * ",
		id, username, password, moniker).Scan(
		&user.ID,
		&user.Username,
		&user.Moniker,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}
