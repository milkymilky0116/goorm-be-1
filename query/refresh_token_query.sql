-- name: CreateRefreshToken :one
INSERT INTO refresh_token (refresh_token,user_id,expires_at) VALUES ($1, $2, $3) RETURNING *;

-- name: GetTokenById :one
SELECT * FROM refresh_token WHERE id = $1;


