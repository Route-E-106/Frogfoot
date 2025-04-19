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

-- name: ReturnGasIncomeHistory :many
SELECT income, change_timestamp FROM gas_income_history
WHERE user_id = ?;

-- name: ReturnMetalIncomeHistory :many
SELECT income, change_timestamp FROM metal_income_history
WHERE user_id = ?;

-- name: UpdateGasIncomeHistory :exec
INSERT INTO gas_income_history(
    income, user_id, change_timestamp
) VALUES (
?, ?, ?
);

-- name: UpdateMetalIncomeHistory :exec
INSERT INTO metal_income_history(
    income, user_id, change_timestamp
) VALUES (
?, ?, ?
);
