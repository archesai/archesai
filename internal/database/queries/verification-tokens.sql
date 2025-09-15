-- name: CreateVerificationToken :one
INSERT INTO
  verification_token (id, identifier, value, expires_at)
VALUES
  ($1, $2, $3, $4)
RETURNING
  *;

-- name: GetVerificationToken :one
SELECT
  *
FROM
  verification_token
WHERE
  id = $1
LIMIT
  1;

-- name: GetVerificationTokenByValue :one
SELECT
  *
FROM
  verification_token
WHERE
  identifier = $1
  AND value = $2
LIMIT
  1;

-- name: ListVerificationTokens :many
SELECT
  *
FROM
  verification_token
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: ListVerificationTokensByIdentifier :many
SELECT
  *
FROM
  verification_token
WHERE
  identifier = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateVerificationToken :one
UPDATE verification_token
SET
  value = COALESCE(sqlc.narg (value), value),
  expires_at = COALESCE(sqlc.narg (expires_at), expires_at),
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  *;

-- name: DeleteVerificationToken :exec
DELETE FROM verification_token
WHERE
  id = $1;

-- name: DeleteVerificationTokensByIdentifier :exec
DELETE FROM verification_token
WHERE
  identifier = $1;
