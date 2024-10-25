package db

import (
	"context"
	"github.com/Darkhackit/simplebank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateRandomAccount(t *testing.T) Account {
	owner := pgtype.Text{
		String: util.RandomOwner(),
		Valid:  true,
	}
	currency := pgtype.Text{
		String: util.RandomCurrency(),
		Valid:  true,
	}
	balance := pgtype.Int8{
		Int64: util.RandomMoney(),
		Valid: true,
	}
	arg := CreateAccountParams{
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	}
	result, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, arg.Owner, result.Owner)
	require.Equal(t, arg.Balance, result.Balance)
	require.Equal(t, arg.Currency, result.Currency)

	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)
	return result
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}
func TestGetAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	ctx := context.Background()
	acc, err := testQueries.GetAccount(ctx, account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, acc.ID, account.ID)
	require.Equal(t, acc.Balance, account.Balance)
	require.Equal(t, acc.Currency, account.Currency)
	require.Equal(t, acc.Owner, acc.Owner)
}

func TestUpdateAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	ctx := context.Background()

	balance := pgtype.Int8{
		Int64: util.RandomMoney(),
		Valid: true,
	}

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: balance,
	}
	result, err := testQueries.UpdateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, arg.ID, result.ID)
	require.Equal(t, arg.Balance, result.Balance)

}
func TestDeleteAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	ctx := context.Background()
	err := testQueries.DeleteAccount(ctx, account.ID)
	require.NoError(t, err)

	acc, err := testQueries.GetAccount(ctx, account.ID)
	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
	require.Empty(t, acc)
}
