-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUserName :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, username, created_at FROM users
ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (
  username, password, created_at
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: ReturnIncomeHistory :many
SELECT resource_name, income, change_timestamp FROM income_history
WHERE user_id = ?;

-- name: UpdateIncomeHistory :exec
INSERT INTO income_history(
    resource_name, income, user_id, change_timestamp
) VALUES (
?, ?, ?, ?
)

