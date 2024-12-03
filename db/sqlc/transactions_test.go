package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bakare-dev/simple-bank-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransaction(t *testing.T) Transaction {
	arg := CreateUserParams{
		Name:        "Test User",
		Email:       util.RandomEmail(),
		Password:    "TestPassword",
		PhoneNumber: sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotZero(t, user.ID)

	accountArg := CreateAccountParams{
		AccountNumber: util.RandomNumberString(12),
		UserID:        user.ID,
		Pin:           util.RandomNumberString(4),
		Type:          AccountTypeSavings,
	}

	account, err := testQueries.CreateAccount(context.Background(), accountArg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)

	transactionArg := CreateTransactionParams{
		Amount:      "2000.00",
		AccountID:   account.ID,
		Type:        TransactionTypeDebit,
		Description: sql.NullString{String: util.RandomString(10), Valid: true},
		Status:      TransactionStatusSuccessful,
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), transactionArg)

	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, transactionArg.Amount, transaction.Amount)
	require.Equal(t, transactionArg.Description, transaction.Description)
	require.Equal(t, transactionArg.Type, transaction.Type)

	require.NotZero(t, transaction.TransactionDate)
	require.NotZero(t, transaction.ID)

	return transaction
}

func TestCreateTransaction(t *testing.T) {
	createRandomTransaction(t)
}

func TestGetTransaction(t *testing.T) {
	transaction1 := createRandomTransaction(t)

	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transaction2)

	require.Equal(t, transaction1.ID, transaction2.ID)
	require.Equal(t, transaction1.Amount, transaction2.Amount)
	require.Equal(t, transaction1.Description, transaction2.Description)
	require.Equal(t, transaction1.Type, transaction2.Type)

	require.WithinDuration(t, transaction1.TransactionDate, transaction2.TransactionDate, time.Second)
}

func TestUpdateTransaction(t *testing.T) {
	transaction1 := createRandomTransaction(t)

	arg := UpdateTransactionParams{
		ID:     transaction1.ID,
		Status: TransactionStatusProcessing,
	}

	transaction2, err := testQueries.UpdateTransaction(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transaction2)

	require.Equal(t, transaction1.ID, transaction2.ID)
	require.Equal(t, arg.Status, transaction2.Status)
}

func TestGetTransactions(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransaction(t)
	}

	arg := ListTransactionsParams{
		Limit:  5,
		Offset: 5,
	}

	transactions, err := testQueries.ListTransactions(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transactions, 5)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
		require.NotZero(t, transaction.ID)
		require.NotZero(t, transaction.Amount)
		require.NotZero(t, transaction.Type)
	}
}

func TestDeleteTransaction(t *testing.T) {
	transaction1 := createRandomTransaction(t)

	err := testQueries.DeleteTransaction(context.Background(), transaction1.ID)

	require.NoError(t, err)

	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transaction2)
}

func TestGetAccountTransaction(t *testing.T) {
	transaction1 := createRandomTransaction(t)

	transaction, err := testQueries.GetUserAccountTransaction(context.Background(), transaction1.AccountID)

	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, transaction1.ID, transaction.ID)
	require.Equal(t, transaction1.Type, transaction.Type)
	require.Equal(t, transaction1.Amount, transaction.Amount)
	require.Equal(t, transaction1.AccountID, transaction.AccountID)
	require.WithinDuration(t, transaction1.TransactionDate, transaction.TransactionDate, time.Second)
}
