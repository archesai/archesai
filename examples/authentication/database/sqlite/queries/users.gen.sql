

-- name: CreateUser :one
INSERT INTO
  "user" (id, email, email_verified, image, name)
VALUES
  (
    $1,
    sqlc.arg('email'),
    sqlc.arg('email_verified'),
    sqlc.narg('image'),
    sqlc.arg('name')
  )
RETURNING
  *;

-- name: GetUser :one
SELECT
  *
FROM
  "user"
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListUsers :many
SELECT
  *
FROM
  "user"
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountUsers :one
SELECT
  COUNT(*)
FROM
  "user";

-- name: UpdateUser :one
UPDATE "user"
SET
  email = COALESCE(sqlc.narg('email'), email),
  email_verified = COALESCE(sqlc.narg('email_verified'), email_verified),
  image = COALESCE(sqlc.narg('image'), image),
  name = COALESCE(sqlc.narg('name'), name)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE
  id = sqlc.arg('id');

-- name: GetUserByEmail :one
SELECT
  *
FROM
  "user"
WHERE
  email = $1
LIMIT
  1;
-- name: GetUserBySessionID :one
SELECT u.*
FROM "user" u
JOIN "session" s ON u.id = s.user_id
WHERE s.id = sqlc.arg(session_id)
LIMIT 1;