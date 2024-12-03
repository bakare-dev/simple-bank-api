-- name: CreateTransaction :one
INSERT INTO "Transactions" (
  account_id, amount, description, status, type
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetTransaction :one
SELECT * 
FROM "Transactions"
WHERE id = $1;

-- name: GetUserAccountTransaction :one
SELECT * 
FROM "Transactions"
WHERE account_id = $1;

-- name: ListTransactions :many
SELECT * 
FROM "Transactions"
ORDER BY transaction_date DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTransaction :one
UPDATE "Transactions"
  set status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM "Transactions"
WHERE id = $1;

-- name: GetAccountBalance :one
SELECT 
  COALESCE(SUM(CASE WHEN type = 'Credit' THEN amount ELSE 0 END), 0) - 
  COALESCE(SUM(CASE WHEN type = 'Debit' THEN amount ELSE 0 END), 0) AS balance
FROM "Transactions"
WHERE account_id = $1;