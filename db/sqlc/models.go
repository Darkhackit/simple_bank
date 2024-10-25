// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        int64            `json:"id"`
	Owner     pgtype.Text      `json:"owner"`
	Balance   pgtype.Int8      `json:"balance"`
	Currency  pgtype.Text      `json:"currency"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type Entries struct {
	ID        int64            `json:"id"`
	AccountID pgtype.Int8      `json:"account_id"`
	Amount    pgtype.Int8      `json:"amount"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type Transfers struct {
	ID            int64            `json:"id"`
	FromAccountID pgtype.Int8      `json:"from_account_id"`
	ToAccountID   pgtype.Int8      `json:"to_account_id"`
	Amount        pgtype.Int8      `json:"amount"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
}
