package user

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *service {
	return &service{db}
}

type User struct {
	Name     string
	Password string
}

func (s *service) AddUser(u User) (string, error) {
	var id string
	q := "INSERT INTO users (username, password) VALUES ($1,$2) RETURNING id"

	stmt, err := s.db.Prepare(q)

	if err != nil {
		return "", fmt.Errorf("creating prepared statement: %w", err)
	}

	err = stmt.QueryRow(u.Name, u.Password).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("query executing: %w", err)
	}

	return id, nil
}
