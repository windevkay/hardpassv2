package mocks

import (
	"time"

	"github.com/windevkay/hardpassv2/internal/entities"
)

var mockPassword = &entities.Password{
	ID:            1,
	App:           "test",
	Password:      "test",
	KeyIdentifier: "test",
	Created_At:    time.Now(),
	Updated_At:    time.Now(),
}

type PasswordEntity struct{}

func (m *PasswordEntity) Insert(appIdentifier string) (int, error) {
	return 1, nil
}

func (m *PasswordEntity) Get(id int) (*entities.Password, error) {
	switch id {
	case 1:
		return mockPassword, nil
	default:
		return nil, entities.ErrNoRecord
	}
}

func (m *PasswordEntity) AllPasswords() ([]*entities.Password, error) {
	return []*entities.Password{mockPassword}, nil
}