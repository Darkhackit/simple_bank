package db

import (
	"context"
	"github.com/Darkhackit/simplebank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateTransfer(t *testing.T) Transfers {
	acct1 := CreateRandomAccount(t)
	acct2 := CreateRandomAccount(t)

	amount := pgtype.Int8{
		Int64: util.RandomMoney(),
		Valid: true,
	}
	id1 := pgtype.Int8{
		Int64: acct1.ID,
		Valid: true,
	}
	id2 := pgtype.Int8{
		Int64: acct2.ID,
		Valid: true,
	}
	arg := CreateTransferParams{
		Amount:        amount,
		FromAccountID: id1,
		ToAccountID:   id2,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateTransfer(t)
}
func TestGetTransfer(t *testing.T) {
	transfer := CreateTransfer(t)
	ctx := context.Background()

	getTransfer, err := testQueries.GetTransfer(ctx, transfer.ID)
	require.NoError(t, err)
	require.Equal(t, transfer, getTransfer)
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := CreateTransfer(t)
	amount := pgtype.Int8{
		Int64: util.RandomMoney(),
		Valid: true,
	}
	arg := UpdateTransferParams{
		Amount: amount,
		ID:     transfer1.ID,
	}
	updatedTrans, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTrans)
	require.Equal(t, amount, updatedTrans.Amount)
	require.Equal(t, transfer1.ID, updatedTrans.ID)

}
func TestDeleteTransfer(t *testing.T) {
	transfer := CreateTransfer(t)
	ctx := context.Background()

	trans := testQueries.DeleteTransfer(ctx, transfer.ID)
	require.NoError(t, trans)

	_, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
}
