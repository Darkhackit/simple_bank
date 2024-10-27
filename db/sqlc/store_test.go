package db

import (
	"context"
	"fmt"
	"github.com/Darkhackit/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB, testQueries)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Print(">> Before", account1.Balance, account2.Balance)
	amount := util.RandomMoney()

	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(arg)

			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.Q.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)

		_, err = store.Q.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)

		_, err = store.Q.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)

		fmt.Println(">> tx \n", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance.Int64 - fromAccount.Balance.Int64
		//diff2 := toAccount.Balance.Int64 - account2.Balance.Int64

		//require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		existed[k] = true

	}

	updatedAccount1, err := store.Q.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	//updatedAccount2, err := store.q.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Print(">> Updated", updatedAccount1.Balance, updatedAccount1.Balance)

	//require.Equal(t, account1.Balance.Int64-int64(n)*amount, updatedAccount1.Balance)
	//require.Equal(t, account2.Balance.Int64+int64(n)*amount, updatedAccount2.Balance)

}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB, testQueries)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Print(">> Before", account1.Balance, account2.Balance)
	amount := util.RandomMoney()

	n := 10
	errs := make(chan error)
	//results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		arg := TransferTxParams{
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
		}
		go func() {
			_, err := store.TransferTx(arg)

			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := store.Q.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	//updatedAccount2, err := store.q.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Print(">> Updated", updatedAccount1.Balance, updatedAccount1.Balance)

	//require.Equal(t, account1.Balance.Int64-int64(n)*amount, updatedAccount1.Balance)
	//require.Equal(t, account2.Balance.Int64+int64(n)*amount, updatedAccount2.Balance)

}
