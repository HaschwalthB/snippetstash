package models

import (
  "database/sql"
  "time"
)


type Snippet struct {
  ID int
  Title string
  Content string
  Created time.Time
  Expires time.Time
}


type SnippetModel struct {
  DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, Expires int) (int, error) {
  return 0, nil
}

func (m *Snippet) Get(id int) (*Snippet, error) {
  return nil, nil
}

func (m *Snippet) Latest() ([]*Snippet, error) {
  return nil, nil
}


