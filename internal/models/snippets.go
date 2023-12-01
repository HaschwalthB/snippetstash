package models

import (
	"database/sql"
	"errors"
	"time"
)

// make a struct to hold the data for a single snippet
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// this is a connection pool
type SnippetModelDB struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModelDB) Insert(title string, content string, expires int) (int, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()
	stmt := `INSERT INTO snippets (title, content, created, expires)
  VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := tx.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	err = tx.Commit()
	return int(id), err
}

// get a specific snippet based on its id on database
func (m *SnippetModelDB) Get(id int) (*Snippet, error) {
	// make a sql statement to retrieve the data
	stmt := `SELECT id, title, content, created, expires FROM snippets
  WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// use QueryRow() to execute the statement and store the result in a new row
	row := m.DB.QueryRow(stmt, id)

	// initialize a pointer to a new zeroed Snippet struct, this will be used to hold the data
	s := &Snippet{}
	// row.Scan() will copy the values from each field in sql.Row to the corresponding field in the Snippet struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// use sentinel error to compare the error returned by the query
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModelDB) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER by created DESC LIMIT 10`
	m.DB.Begin()

	rows, err := m.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// close connection before Latest() returns
	// initialize an empty slice to hold the returned snippets
	defer rows.Close()
	snippets := []*Snippet{}
	// use rows.Next() with for loop to iterate each row
	for rows.Next() {
		// initialize a pointer to a new zeroed Snippet struct
		// copy the values from each field in sql.Row to the corresponding field in the Snippet struct
		// and append it to the slice of Snippets
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
