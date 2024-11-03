package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Queries
}

type SQLStore struct {
	Q  *Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool, q *Queries) SQLStore {
	return SQLStore{
		db: db,
		Q:  q,
	}
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Account   `json:"from_account"`
	ToAccount   Account   `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

func (s *SQLStore) TransferTx(arg TransferTxParams) (TransferTxResult, error) {
	tx, err := s.db.Begin(context.Background())
	var result TransferTxResult
	if err != nil {
		return TransferTxResult{}, err
	}
	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil {
			return
		}
	}()
	qtx := s.Q.WithTx(tx)
	fromID := pgtype.Int8{
		Int64: arg.FromAccountID,
		Valid: true,
	}
	toID := pgtype.Int8{
		Int64: arg.ToAccountID,
		Valid: true,
	}
	amount := pgtype.Int8{
		Int64: arg.Amount,
		Valid: true,
	}
	negAmount := pgtype.Int8{
		Int64: -arg.Amount,
		Valid: true,
	}
	result.Transfer, err = qtx.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        amount,
	})
	if err != nil {
		return TransferTxResult{}, err
	}
	result.FromEntry, err = qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: fromID,
		Amount:    negAmount,
	})
	if err != nil {
		return TransferTxResult{}, err
	}
	result.ToEntry, err = qtx.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: fromID,
		Amount:    amount,
	})
	if err != nil {
		return TransferTxResult{}, err
	}

	//account1, err := qtx.GetAccountForUpdate(context.Background(), arg.FromAccountID)
	//if err != nil {
	//	return TransferTxResult{}, err
	//}
	//b := pgtype.Int8{
	//	Int64: account1.Balance.Int64 - arg.Amount,
	//	Valid: true,
	//}
	if arg.FromAccountID < arg.ToAccountID {
		//b := pgtype.Int8{
		//	Int64: -arg.Amount,
		//	Valid: true,
		//}
		//result.FromAccount, err = qtx.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		//	ID:     arg.FromAccountID,
		//	Amount: b,
		//})
		//if err != nil {
		//	return TransferTxResult{}, err
		//}
		////account2, err := qtx.Get(context.Background(), arg.FromAccountID)
		////if err != nil {
		////	return TransferTxResult{}, err
		////}
		//bp := pgtype.Int8{
		//	Int64: arg.Amount,
		//	Valid: true,
		//}
		//result.ToAccount, err = qtx.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		//	ID:     arg.ToAccountID,
		//	Amount: bp,
		//})
		//if err != nil {
		//	return TransferTxResult{}, err
		//}
		result.FromAccount, result.ToAccount, err = addMoney(context.Background(), qtx, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		if err != nil {
			return TransferTxResult{}, err
		}
	} else {
		//bp := pgtype.Int8{
		//	Int64: arg.Amount,
		//	Valid: true,
		//}
		//result.ToAccount, err = qtx.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		//	ID:     arg.ToAccountID,
		//	Amount: bp,
		//})
		//if err != nil {
		//	return TransferTxResult{}, err
		//}
		//b := pgtype.Int8{
		//	Int64: -arg.Amount,
		//	Valid: true,
		//}
		//result.FromAccount, err = qtx.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		//	ID:     arg.FromAccountID,
		//	Amount: b,
		//})
		//if err != nil {
		//	return TransferTxResult{}, err
		//}
		result.ToAccount, result.FromAccount, err = addMoney(context.Background(), qtx, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

	}
	err = tx.Commit(context.Background())
	if err != nil {
		return TransferTxResult{}, err
	}
	return result, nil
}

func addMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	a1 := pgtype.Int8{
		Int64: amount1,
		Valid: true,
	}
	a2 := pgtype.Int8{
		Int64: amount2,
		Valid: true,
	}
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: a1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: a2,
	})
	if err != nil {
		return
	}
	return account1, account2, nil
}
