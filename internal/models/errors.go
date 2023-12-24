package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials is used when a user tries to login with an incorrect password
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail is used when a user tries to signup with an email that is already in use
	ErrDuplicateEmail = errors.New("models : duplicate email")
)
