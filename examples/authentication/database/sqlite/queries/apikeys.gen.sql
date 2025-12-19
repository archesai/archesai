

-- name: CreateAPIKey :one
INSERT INTO
  api_key (id, expires_at, key_hash, last_used_at, name, organization_id, prefix, rate_limit, scopes, user_id)
VALUES
  (
    $1,
    sqlc.narg('expires_at'),
    sqlc.arg('key_hash'),
    sqlc.narg('last_used_at'),
    sqlc.narg('name'),
    sqlc.arg('organization_id'),
    sqlc.narg('prefix'),
    sqlc.arg('rate_limit'),
    sqlc.arg('scopes'),
    sqlc.arg('user_id')
  )
RETURNING
  *;

-- name: GetAPIKey :one
SELECT
  *
FROM
  api_key
WHERE
  id = sqlc.arg('id')
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
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountAPIKeys :one
SELECT
  COUNT(*)
FROM
  api_key;

-- name: UpdateAPIKey :one
UPDATE api_key
SET
  expires_at = COALESCE(sqlc.narg('expires_at'), expires_at),
  key_hash = COALESCE(sqlc.narg('key_hash'), key_hash),
  last_used_at = COALESCE(sqlc.narg('last_used_at'), last_used_at),
  name = COALESCE(sqlc.narg('name'), name),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id),
  prefix = COALESCE(sqlc.narg('prefix'), prefix),
  rate_limit = COALESCE(sqlc.narg('rate_limit'), rate_limit),
  scopes = COALESCE(sqlc.narg('scopes'), scopes),
  user_id = COALESCE(sqlc.narg('user_id'), user_id)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteAPIKey :exec
DELETE FROM api_key
WHERE
  id = sqlc.arg('id');