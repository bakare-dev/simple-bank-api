package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bakare-dev/simple-bank-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
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

	require.Equal(t, accountArg.AccountNumber, account.AccountNumber)
	require.Equal(t, accountArg.Pin, account.Pin)
	require.Equal(t, accountArg.UserID, account.UserID)

	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.UpdatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.AccountNumber, account2.AccountNumber)
	require.Equal(t, account1.Pin, account2.Pin)
	require.Equal(t, account1.Type, account2.Type)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:  account1.ID,
		Pin: util.RandomNumberString(4),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Pin, account2.Pin)
}

func TestGetAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.NotZero(t, account.ID)
		require.NotZero(t, account.UserID)
		require.NotZero(t, account.AccountNumber)
		require.NotZero(t, account.CreatedAt)
	}
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestGetUserAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account, err := testQueries.GetUserAccount(context.Background(), account1.UserID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account1.ID, account.ID)
	require.Equal(t, account1.AccountNumber, account.AccountNumber)
	require.Equal(t, account1.Pin, account.Pin)
	require.Equal(t, account1.Type, account.Type)
	require.Equal(t, account1.UserID, account.UserID)
	require.WithinDuration(t, account1.CreatedAt, account.CreatedAt, time.Second)
}
