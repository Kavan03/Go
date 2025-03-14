-- -- name: CreateUser :one → This is a special sqlc comment that defines a query named CreateUser.
-- :one → Specifies that the query returns exactly one row (a single User struct in Go).
-- RETURNING * → Returns the newly inserted row.

-- name: CreateUser :one
INSERT INTO users (id,created_at,updated_at,name,api_key)
VALUES ($1,$2,$3,$4,
encode(sha256(random()::text::bytea),'hex')
)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users where api_key = $1;