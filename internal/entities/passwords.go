package entities

import (
	"database/sql"
	"errors"
	"time"
)

type Password struct {
	ID         int
	App        string
	Password   string
	Created_At time.Time
	Updated_At time.Time
}

type PasswordEntity struct {
	DB *sql.DB
}

func (p *PasswordEntity) Insert(app string) (int, error) {
	// todo: generate password
	//	- https://pkg.go.dev/crypto/rand#Read
	// todo: hash password
	//	- https://pkg.go.dev/golang.org/x/crypto/bcrypt#GenerateFromPassword
	password := "place_holder"
	stmt := `INSERT INTO passwords (app, password) VALUES (?, ?)`
	result, err := p.DB.Exec(stmt, app, password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (p *PasswordEntity) Get(id int) (*Password, error) {
	stmt := `SELECT id, app, password, created_at, updated_at FROM passwords WHERE id = ?`
	row := p.DB.QueryRow(stmt, id)

	password := &Password{}
	err := row.Scan(&password.ID, &password.App, &password.Password, &password.Created_At, &password.Updated_At)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// todo: decrypt hashed password
	//	- https://pkg.go.dev/golang.org/x/crypto/bcrypt#CompareHashAndPassword
	return password, nil
}

func (p *PasswordEntity) AllPasswords() ([]*Password, error) {
	stmt := `SELECT id, app, password, created_at, updated_at 
	FROM passwords
	WHERE deleted_at IS NULL`

	rows, err := p.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	passwords := []*Password{}
	for rows.Next() {
		password := &Password{}
		err := rows.Scan(&password.ID, &password.App, &password.Password, &password.Created_At, &password.Updated_At)
		if err != nil {
			return nil, err
		}
		passwords = append(passwords, password)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return passwords, nil
}