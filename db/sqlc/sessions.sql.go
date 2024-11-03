// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: sessions.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions(
    id, username, refresh_token, user_agent, client_ip
) VALUES ($1 ,$2 ,$3,$4,$5) RETURNING id, username, refresh_token, user_agent, client_ip, is_blocked, created_at, expires_at
`

type CreateSessionParams struct {
	ID           pgtype.UUID `json:"id"`
	Username     string      `json:"username"`
	RefreshToken string      `json:"refresh_token"`
	UserAgent    string      `json:"user_agent"`
	ClientIp     string      `json:"client_ip"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Sessions, error) {
	row := q.db.QueryRow(ctx, createSession,
		arg.ID,
		arg.Username,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
	)
	var i Sessions
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT id, username, refresh_token, user_agent, client_ip, is_blocked, created_at, expires_at FROM sessions WHERE id = $1
`

func (q *Queries) GetSession(ctx context.Context, id pgtype.UUID) (Sessions, error) {
	row := q.db.QueryRow(ctx, getSession, id)
	var i Sessions
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}
