package store

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type SqlDb struct {
	db *sql.DB
}

func NewSqlDB(url string) *SqlDb {
	db, err := sql.Open("sqlite", url)

	if err != nil {
		log.Fatal("Failed open db connection", err)
	}

	return &SqlDb{
		db: db,
	}
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
