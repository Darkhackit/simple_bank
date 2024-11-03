// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
     email, username, active, name, password
) VALUES ($1 ,$2 , $3 ,$4,$5) RETURNING id, email, username, active, name, password, password_changed_at, created_at
`

type CreateUserParams struct {
	Email    string      `json:"email"`
	Username string      `json:"username"`
	Active   pgtype.Bool `json:"active"`
	Name     pgtype.Text `json:"name"`
	Password string      `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (Users, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.Username,
		arg.Active,
		arg.Name,
		arg.Password,
	)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Active,
		&i.Name,
		&i.Password,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, email, username, active, name, password, password_changed_at, created_at FROM users WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (Users, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Active,
		&i.Name,
		&i.Password,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
