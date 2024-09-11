-- name: GetUser :one
SELECT * FROM "user" WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (email, password, role) VALUES ($1, $2, $3) RETURNING *;
