package models

import (
	"database/sql"
	"errors"
	"time"
)

// Represents a single snippet in the database
type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// wrapper for sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	// get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	s := &Snippet{}

	// query the database for a snippet with the given ID, then copy the values into the Snippet struct
	err := m.DB.QueryRow(query, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// check if no matching record is found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}

		return nil, err
	}

	return s, nil;
}

// returns the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// defer the closing of the rows until the function returns
	defer rows.Close()

	snippets := []*Snippet{}

	// iterate through the rows with rows.Next() and copy the values from each row into a Snippet struct
	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// append new snippet to the slice of snippets
		snippets = append(snippets, s)
	}

	// check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
