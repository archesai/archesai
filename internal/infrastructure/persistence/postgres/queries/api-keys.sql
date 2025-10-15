-- name: CreateAPIKey :one
INSERT INTO
  api_key (
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

-- name: GetAPIKey :one
SELECT
  *
FROM
  api_key
WHERE
  id = $1
LIMIT
  1;

-- name: GetAPIKeyByKeyHash :one
SELECT
  *
FROM
  api_key
WHERE
  key_hash = $1
LIMIT
  1;

-- name: ListAPIKeys :many
SELECT
  *
FROM
  api_key
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: ListAPIKeysByUser :many
SELECT
  *
FROM
  api_key
WHERE
  user_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateAPIKey :one
UPDATE api_key
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

-- name: DeleteAPIKey :exec
DELETE FROM api_key
WHERE
  id = $1;

-- name: DeleteAPIKeysByUser :exec
DELETE FROM api_key
WHERE
  user_id = $1;

-- name: ListAPIKeysByOrganization :many
SELECT
  *
FROM
  api_key
WHERE
  organization_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_key
SET
  last_used_at = CURRENT_TIMESTAMP
WHERE
  id = $1;

-- name: DeleteExpiredAPIKeys :exec
DELETE FROM api_key
WHERE
  expires_at IS NOT NULL
  AND expires_at < CURRENT_TIMESTAMP;
