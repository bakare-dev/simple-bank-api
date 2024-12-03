-- name: CreateAccount :one
INSERT INTO "Accounts" (
  user_id, account_number, pin, type
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetAccount :one
SELECT * 
FROM "Accounts"
WHERE id = $1;

-- name: GetUserAccount :one
SELECT * 
FROM "Accounts"
WHERE user_id = $1;

-- name: ListAccounts :many
SELECT * 
FROM "Accounts"
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE "Accounts"
  set pin = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM "Accounts"
WHERE id = $1;