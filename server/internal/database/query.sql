-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY userName;

-- name: CreateUser :one
INSERT INTO users (
  userName, password, created_at
) VALUES (
  ?, ?, ?
)
RETURNING *;
