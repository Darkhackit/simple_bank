package db

import (
	"context"
	"github.com/Darkhackit/simplebank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateEntry(t *testing.T) Entries {
	randomAcc := CreateRandomAccount(t)

	id := pgtype.Int8{
		Int64: randomAcc.ID,
		Valid: true,
	}
	money := pgtype.Int8{
		Int64: util.RandomMoney(),
		Valid: true,
	}
	arg := CreateEntryParams{
		AccountID: id,
		Amount:    money,
	}
	ctx := context.Background()
	entry, err := testQueries.CreateEntry(ctx, arg)
	require.NoError(t, err)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)
	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateEntry(t)
}

func TestGetEntry(t *testing.T) {
	ent := CreateEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), ent.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, ent, entry)
}

func TestUpdateEntry(t *testing.T) {
	ent := CreateEntry(t)
	money := pgtype.Int8{
		Int64: util.RandomMoney(),
		Valid: true,
	}
	arg := UpdateEntryParams{
		ID:     ent.ID,
		Amount: money,
	}
	update, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Amount, update.Amount)
	require.Equal(t, arg.ID, update.ID)
}

func TestDeleteEntry(t *testing.T) {
	ent := CreateEntry(t)
	err := testQueries.DeleteEntry(context.Background(), ent.ID)
	require.NoError(t, err)
	_, err = testQueries.GetEntry(context.Background(), ent.ID)
	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
}
