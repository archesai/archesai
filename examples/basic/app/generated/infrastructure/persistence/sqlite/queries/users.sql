

-- name: CreateUser :one
INSERT INTO
  "user" (id, email, name)
VALUES
  (
    $1,
    sqlc.arg('email'),
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
  name = COALESCE(sqlc.narg('name'), name)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE
  id = sqlc.arg('id');