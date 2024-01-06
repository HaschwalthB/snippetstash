package models

import (
	"database/sql"
	"time"
)

type users struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// open connection pool to the database
type UserModelDB struct {
	*sql.DB
}

// Insert adds a new record to the users table.
func (m *UserModelDB) Insert(name, email, password string) error {
	return nil
}

// Authenticate verifies whether a user exists with the provided email address and password.
func (m *UserModelDB) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// check if a user already exists
func (m *UserModelDB) Exists(id int) (bool, error) {
	return false, nil
}
