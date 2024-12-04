package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        string    `json:"amount"`
}

type TransferTxResult struct {
	DebitTransaction  Transaction    `json:"debittransaction"`
	CreditTransaction Transaction    `json:"credittransaction"`
	FromAccount       AccountBalance `json:"from_account"`
	ToAccount         AccountBalance `json:"to_account"`
}

type AccountBalance struct {
	Balance string  `json:"balance"`
	Account Account `json:"account"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var fromAccount, toAccount uuid.UUID
		fromAccount = arg.ToAccountID
		toAccount = arg.FromAccountID

		_, err := q.GetAccount(ctx, fromAccount)
		if err != nil {
			return fmt.Errorf("failed to lock from account: %w", err)
		}

		_, err = q.GetAccount(ctx, toAccount)
		if err != nil {
			return fmt.Errorf("failed to lock to account: %w", err)
		}

		fromBalance, err := q.GetAccountBalance(ctx, arg.FromAccountID)
		if err != nil {
			return fmt.Errorf("failed to get balance for from account: %w", err)
		}

		if fromBalance < arg.Amount {
			return fmt.Errorf("insufficient funds in from account")
		}

		result.DebitTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			Amount:      arg.Amount,
			Status:      TransactionStatusSuccessful,
			Description: sql.NullString{String: "transfer", Valid: true},
			Type:        TransactionTypeDebit,
			AccountID:   arg.FromAccountID,
		})
		if err != nil {
			return fmt.Errorf("failed to create debit transaction: %w", err)
		}

		result.CreditTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			Amount:      arg.Amount,
			Status:      TransactionStatusSuccessful,
			Description: sql.NullString{String: "transfer", Valid: true},
			Type:        TransactionTypeCredit,
			AccountID:   arg.ToAccountID,
		})
		if err != nil {
			return fmt.Errorf("failed to create credit transaction: %w", err)
		}

		result.FromAccount.Account, err = q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return fmt.Errorf("failed to fetch from account details: %w", err)
		}
		result.FromAccount.Balance, err = q.GetAccountBalance(ctx, arg.FromAccountID)
		if err != nil {
			return fmt.Errorf("failed to update from account balance: %w", err)
		}

		result.ToAccount.Account, err = q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return fmt.Errorf("failed to fetch to account details: %w", err)
		}
		result.ToAccount.Balance, err = q.GetAccountBalance(ctx, arg.ToAccountID)
		if err != nil {
			return fmt.Errorf("failed to update to account balance: %w", err)
		}

		return nil
	})

	return result, err
}
