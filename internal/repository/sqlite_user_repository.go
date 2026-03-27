package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"testsystem/internal/model"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) (*SQLiteUserRepository, error) {
	repo := &SQLiteUserRepository{db: db}
	if err := repo.migrate(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *SQLiteUserRepository) migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL
)`)
	return err
}

func (r *SQLiteUserRepository) Create(ctx context.Context, name string) (model.User, error) {
	if strings.TrimSpace(name) == "" {
		return model.User{}, errors.New("name cannot be empty")
	}

	res, err := r.db.ExecContext(ctx, `INSERT INTO users(name) VALUES(?)`, name)
	if err != nil {
		return model.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.User{}, err
	}
	return model.User{ID: id, Name: name}, nil
}

func (r *SQLiteUserRepository) List(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
