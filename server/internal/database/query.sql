-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUserName :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: GetUserMetalExtractorLevel :one
SELECT metal_extractor_lvl FROM users
WHERE id = ?;
-- name: GetUserGasExtractorLevel :one
SELECT gas_extractor_lvl FROM users
WHERE id = ?;

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

-- name: UpdateUserMetalExtractorLevel :exec
UPDATE users SET metal_extractor_lvl = 1 + users.metal_extractor_lvl
WHERE id = ?;

-- name: UpdateUserGasExtractorLevel :exec
UPDATE users SET gas_extractor_lvl = 1 + users.gas_extractor_lvl
WHERE id = ?;

-- name: GetUserExpenses :one
SELECT total_gas_expenses, total_metal_expenses FROM users
WHERE id = ? LIMIT 1;

-- name: UpdateTotalExpenses :exec
UPDATE users SET total_gas_expenses = ? + users.total_gas_expenses, total_metal_expenses = ? + users.total_metal_expenses
WHERE id = ?;

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
