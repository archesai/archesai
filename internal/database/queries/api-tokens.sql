-- name: CreateApiToken :one
INSERT INTO
  api_token (
    id,
    user_id,
    organization_id,
    name,
    key_hash,
    prefix,
    scopes,
    rate_limit,
    expires_at
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
  *;

-- name: GetApiToken :one
SELECT
  *
FROM
  api_token
WHERE
  id = $1
LIMIT
  1;

-- name: GetApiTokenByKeyHash :one
SELECT
  *
FROM
  api_token
WHERE
  key_hash = $1
LIMIT
  1;

-- name: ListApiTokens :many
SELECT
  *
FROM
  api_token
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: ListApiTokensByUser :many
SELECT
  *
FROM
  api_token
WHERE
  user_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateApiToken :one
UPDATE api_token
SET
  name = COALESCE(sqlc.narg (name), name),
  scopes = COALESCE(sqlc.narg (scopes), scopes),
  rate_limit = COALESCE(sqlc.narg (rate_limit), rate_limit),
  expires_at = COALESCE(sqlc.narg (expires_at), expires_at),
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = sqlc.arg (id)
RETURNING
  *;

-- name: DeleteApiToken :exec
DELETE FROM api_token
WHERE
  id = $1;

-- name: DeleteApiTokensByUser :exec
DELETE FROM api_token
WHERE
  user_id = $1;

-- name: ListApiTokensByOrganization :many
SELECT
  *
FROM
  api_token
WHERE
  organization_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateApiTokenLastUsed :exec
UPDATE api_token
SET
  last_used_at = CURRENT_TIMESTAMP
WHERE
  id = $1;

-- name: DeleteExpiredApiTokens :exec
DELETE FROM api_token
WHERE
  expires_at IS NOT NULL
  AND expires_at < CURRENT_TIMESTAMP;
