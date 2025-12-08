-- name: CreateUser :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;


-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users 
SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: FindOrCreateUser :one
INSERT INTO users (email, provider, phone, fullname, password_hash, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (email)
DO UPDATE SET
    phone = COALESCE(NULLIF(EXCLUDED.phone, ''), users.phone),
    provider = EXCLUDED.provider,
    fullname = COALESCE(NULLIF(EXCLUDED.fullname, ''), users.fullname),
    password_hash = EXCLUDED.password_hash,
    updated_at = NOW()
RETURNING id, email, fullname, password_hash, phone, role, provider;
