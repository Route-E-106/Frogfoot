-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUserName :one
SELECT * FROM users
WHERE userName = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, userName, created_at FROM users
ORDER BY userName;

-- name: CreateUser :one
INSERT INTO users (
  userName, password, created_at
) VALUES (
  ?, ?, ?
)
RETURNING *;
