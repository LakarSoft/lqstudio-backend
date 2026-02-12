-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1::varchar LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetAdminByEmail :one
SELECT * FROM users
WHERE email = $1 AND role = 'ADMIN' LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    id, name, email, phone, password_hash, role
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    name = $2,
    phone = $3,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;
