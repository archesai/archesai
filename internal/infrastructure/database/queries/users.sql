-- name: CreateUser :one
INSERT INTO "user" (
    email,
    name,
    email_verified,
    image
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM "user"
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE "user"
SET 
    name = COALESCE(sqlc.narg(name), name),
    email = COALESCE(sqlc.narg(email), email),
    email_verified = COALESCE(sqlc.narg(email_verified), email_verified),
    image = COALESCE(sqlc.narg(image), image)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;