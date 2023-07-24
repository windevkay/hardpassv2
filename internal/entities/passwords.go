package entities

import (
	"database/sql"
	"errors"
	"time"

	"github.com/windevkay/hardpassv2/internal/azure"
)

type PasswordEntityInterface interface {
	Insert(appIdentifier string) (int, error)
	Get(id int) (*Password, error)
	AllPasswords() ([]*Password, error)
}

type Password struct {
	ID            int
	App           string
	Password      string
	KeyIdentifier string
	Created_At    time.Time
	Updated_At    time.Time
}

type PasswordEntity struct {
	DB          *sql.DB
	AzureClient *azure.AzureClient
}

func (p *PasswordEntity) Insert(appIdentifier string) (int, error) {
	password, err := azure.GenPassword(p.AzureClient, appIdentifier, nil)
	if err != nil {
		return 0, errors.New("error generating password")
	}
	text := password.Text
	keyIdentifier := password.KeyIdentifier
	// begin transaction
	tx, err := p.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	stmt := `INSERT INTO passwords (app, password, keyIdentifier) VALUES (?, ?, ?)`
	result, err := p.DB.Exec(stmt, appIdentifier, text, keyIdentifier)

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
	// begin transaction
	tx, err := p.DB.Begin()
	if err != nil {
		return &Password{}, err
	}

	defer tx.Rollback()

	stmt := `SELECT id, app, password, keyIdentifier, created_at, updated_at FROM passwords WHERE id = ?`
	row := p.DB.QueryRow(stmt, id)

	password := &Password{}

	err = row.Scan(&password.ID, &password.App, &password.Password, &password.KeyIdentifier, &password.Created_At, &password.Updated_At)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	decryptedPassword, err := azure.DecryptPassword(p.AzureClient, password.KeyIdentifier, password.Password)
	if err != nil {
		return nil, err
	}
	password.Password = decryptedPassword

	return password, nil
}

func (p *PasswordEntity) AllPasswords() ([]*Password, error) {
	stmt := `SELECT id, app, created_at, updated_at 
	FROM passwords
	WHERE deleted_at IS NULL`
	// begin transaction
	tx, err := p.DB.Begin()
	if err != nil {
		return []*Password{}, err
	}

	defer tx.Rollback()

	rows, err := p.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	passwords := []*Password{}
	for rows.Next() {
		password := &Password{}
		err := rows.Scan(&password.ID, &password.App, &password.Created_At, &password.Updated_At)
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
