-- name: CreateUserProfile :one
INSERT INTO user_profile (user_id, is_premium, stripe_id, name, bio)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserProfile :one
SELECT * FROM user_profile 
WHERE user_id = $1;


-- name: UpdateUserProfile :one
UPDATE user_profile 
SET 
    name = COALESCE($2, name),
    bio = COALESCE($3, bio),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;

-- name: UpdatePremiumStatus :one
UPDATE user_profile 
SET is_premium = $2, stripe_id = $3, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;
