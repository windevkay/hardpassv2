package entities

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Email          string
	Name           string
	HashedPassword []byte
	Created        time.Time
}

type UserEntity struct {
	DB *sql.DB
}

func (u *UserEntity) Insert(email, name, password string) (int, error) {
	// begin transaction
	tx, err := u.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	stmt := `INSERT INTO users (email, name, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`
	result, err := u.DB.Exec(stmt, email, name, password)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (u *UserEntity) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (u *UserEntity) Exists(id int) (bool, error) {
	return false, nil
}
