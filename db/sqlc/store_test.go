package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bakare-dev/simple-bank-api/util"
	"github.com/stretchr/testify/require"
)

func createTestTransaction(t *testing.T, account Account, amount string, txType TransactionType) Transaction {
	params := CreateTransactionParams{
		Amount:      amount,
		AccountID:   account.ID,
		Type:        txType,
		Description: sql.NullString{String: util.RandomString(10), Valid: true},
		Status:      TransactionStatusSuccessful,
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)
	require.Equal(t, params.Amount, transaction.Amount)

	return transaction
}

func executeConcurrentTransfers(t *testing.T, store *Store, account1, account2 Account, n int, amount string, errs chan error, results chan TransferTxResult) {
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}
}
func executeDeadlockTransfers(t *testing.T, store *Store, account1, account2 Account, n int, amount string, errs chan error) {
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func(fromID, toID int64) {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromID,
				ToAccountID:   toID,
				Amount:        amount,
			})

			errs <- err
		}(fromAccountID, toAccountID)
	}
}

func verifyTransferResult(t *testing.T, result TransferTxResult, account1, account2 Account, amount string) {
	debitTransaction := result.DebitTransaction
	require.NotEmpty(t, debitTransaction)
	require.Equal(t, account1.ID, debitTransaction.AccountID)
	require.Equal(t, amount, debitTransaction.Amount)
	require.NotZero(t, debitTransaction.ID)
	require.NotZero(t, debitTransaction.TransactionDate)

	creditTransaction := result.CreditTransaction
	require.NotEmpty(t, creditTransaction)
	require.Equal(t, account2.ID, creditTransaction.AccountID)
	require.Equal(t, amount, creditTransaction.Amount)
	require.NotZero(t, creditTransaction.ID)
	require.NotZero(t, creditTransaction.TransactionDate)
}

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	createTestTransaction(t, account1, "10000000.00", TransactionTypeCredit)
	createTestTransaction(t, account1, "10000000.00", TransactionTypeCredit)

	n := 5
	amount := "2000.00"
	errs := make(chan error)
	results := make(chan TransferTxResult)

	executeConcurrentTransfers(t, store, account1, account2, n, amount, errs, results)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		verifyTransferResult(t, result, account1, account2, amount)

		_, err = store.GetTransaction(context.Background(), result.DebitTransaction.ID)
		require.NoError(t, err)

		_, err = store.GetTransaction(context.Background(), result.CreditTransaction.ID)
		require.NoError(t, err)
	}
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	createTestTransaction(t, account1, "10000000.00", TransactionTypeCredit)
	createTestTransaction(t, account1, "10000000.00", TransactionTypeCredit)

	n := 10
	amount := "2000.00"
	errs := make(chan error, n)

	executeDeadlockTransfers(t, store, account1, account2, n, amount, errs)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
}
