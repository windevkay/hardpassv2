package entities

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	secure "github.com/windevkay/hardpassv2/internal/utils"
)

type Password struct {
	ID         int
	App        string
	Password   string
	Ekey       string
	Created_At time.Time
	Updated_At time.Time
}

type PasswordEntity struct {
	sync.RWMutex
	DB *sql.DB
}

func (p *PasswordEntity) Insert(appIdentifier string) (int, error) {
	password, err := secure.GenPassword()
	if err != nil {
		return 0, errors.New("error generating password")
	}
	text := password.Text
	key := password.Key
	// begin transaction
	tx, err := p.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	stmt := `INSERT INTO passwords (app, password, ekey) VALUES (?, ?, ?)`
	p.Lock()
	result, err := p.DB.Exec(stmt, appIdentifier, text, key)
	p.Unlock()

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

	stmt := `SELECT id, app, password, ekey, created_at, updated_at FROM passwords WHERE id = ?`
	p.RLock()
	row := p.DB.QueryRow(stmt, id)
	p.RUnlock()

	password := &Password{}

	err = row.Scan(&password.ID, &password.App, &password.Password, &password.Ekey, &password.Created_At, &password.Updated_At)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	decodedPassword, err := hex.DecodeString(password.Password)
	if err != nil {
		return nil, err
	}
	decodedKey, err := hex.DecodeString(password.Ekey)
	if err != nil {
		return nil, err
	}

	decryptedPassword, err := secure.Decrypt(decodedKey, decodedPassword)
	if err != nil {
		return nil, err
	}
	password.Password = hex.EncodeToString(decryptedPassword)

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

	p.RLock()
	rows, err := p.DB.Query(stmt)
	p.RUnlock()

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
