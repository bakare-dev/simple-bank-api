// Code generated by sqlc. DO NOT EDIT.
// source: transactions.sql

package db

import (
	"context"
	"database/sql"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO "Transactions" (
  account_id, amount, description, status, type
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, account_id, type, status, amount, description, transaction_date
`

type CreateTransactionParams struct {
	AccountID   int64             `json:"account_id"`
	Amount      string            `json:"amount"`
	Description sql.NullString    `json:"description"`
	Status      TransactionStatus `json:"status"`
	Type        TransactionType   `json:"type"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.AccountID,
		arg.Amount,
		arg.Description,
		arg.Status,
		arg.Type,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Type,
		&i.Status,
		&i.Amount,
		&i.Description,
		&i.TransactionDate,
	)
	return i, err
}

const deleteTransaction = `-- name: DeleteTransaction :exec
DELETE FROM "Transactions"
WHERE id = $1
`

func (q *Queries) DeleteTransaction(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransaction, id)
	return err
}

const getAccountBalance = `-- name: GetAccountBalance :one
SELECT 
  CAST(
    COALESCE(SUM(CASE WHEN type = 'Credit' THEN amount ELSE 0 END), 0) - 
    COALESCE(SUM(CASE WHEN type = 'Debit' THEN amount ELSE 0 END), 0)
  AS DECIMAL(20, 2)) AS balance
FROM "Transactions"
WHERE account_id = $1
`

func (q *Queries) GetAccountBalance(ctx context.Context, accountID int64) (string, error) {
	row := q.db.QueryRowContext(ctx, getAccountBalance, accountID)
	var balance string
	err := row.Scan(&balance)
	return balance, err
}

const getTransaction = `-- name: GetTransaction :one
SELECT id, account_id, type, status, amount, description, transaction_date 
FROM "Transactions"
WHERE id = $1
`

func (q *Queries) GetTransaction(ctx context.Context, id int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Type,
		&i.Status,
		&i.Amount,
		&i.Description,
		&i.TransactionDate,
	)
	return i, err
}

const getUserAccountTransaction = `-- name: GetUserAccountTransaction :one
SELECT id, account_id, type, status, amount, description, transaction_date 
FROM "Transactions"
WHERE account_id = $1
`

func (q *Queries) GetUserAccountTransaction(ctx context.Context, accountID int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getUserAccountTransaction, accountID)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Type,
		&i.Status,
		&i.Amount,
		&i.Description,
		&i.TransactionDate,
	)
	return i, err
}

const listTransactions = `-- name: ListTransactions :many
SELECT id, account_id, type, status, amount, description, transaction_date 
FROM "Transactions"
ORDER BY transaction_date DESC
LIMIT $1 OFFSET $2
`

type ListTransactionsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, listTransactions, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Type,
			&i.Status,
			&i.Amount,
			&i.Description,
			&i.TransactionDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTransaction = `-- name: UpdateTransaction :one
UPDATE "Transactions"
  set status = $2
WHERE id = $1
RETURNING id, account_id, type, status, amount, description, transaction_date
`

type UpdateTransactionParams struct {
	ID     int64             `json:"id"`
	Status TransactionStatus `json:"status"`
}

func (q *Queries) UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, updateTransaction, arg.ID, arg.Status)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Type,
		&i.Status,
		&i.Amount,
		&i.Description,
		&i.TransactionDate,
	)
	return i, err
}
