-- name: CreateUser :one
INSERT INTO "Users" (
  name, email, password, phone_number
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * 
FROM "Users"
WHERE id = $1;

-- name: ListUsers :many
SELECT * 
FROM "Users"
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE "Users"
  set password = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "Users"
WHERE id = $1;