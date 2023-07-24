package mocks

import (
	"github.com/windevkay/hardpassv2/internal/entities"
)

type UserEntity struct{}

func (m *UserEntity) Insert(name, email, password string) (int, error) { 
	switch email {
	case "dupe@example.com":
		return 0, entities.ErrDuplicateEmail
	default: 
		return 1, nil
	} 
}

func (m *UserEntity) Authenticate(email, password string) (int, error) { 
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil 
	}

	return 0, entities.ErrInvalidCredentials 
}

func (m *UserEntity) Exists(id int) (bool, error) { 
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	} 
}