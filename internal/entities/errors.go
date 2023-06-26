package entities

import (
	"errors"
)

var ErrNoRecord = errors.New("entities: no matching record found")
var ErrInvalidCredentials = errors.New("entities: invalid credentials")
var ErrDuplicateEmail = errors.New("entities: duplicate email")
