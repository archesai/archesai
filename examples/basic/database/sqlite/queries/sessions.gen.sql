

-- name: CreateSession :one
INSERT INTO
  "session" (id, auth_method, auth_provider, expires_at, ip_address, organization_id, token, user_agent, user_id)
VALUES
  (
    $1,
    sqlc.narg('auth_method'),
    sqlc.narg('auth_provider'),
    sqlc.arg('expires_at'),
    sqlc.narg('ip_address'),
    sqlc.narg('organization_id'),
    sqlc.arg('token'),
    sqlc.narg('user_agent'),
    sqlc.arg('user_id')
  )
RETURNING
  *;

-- name: GetSession :one
SELECT
  *
FROM
  "session"
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListSessions :many
SELECT
  *
FROM
  "session"
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountSessions :one
SELECT
  COUNT(*)
FROM
  "session";

-- name: UpdateSession :one
UPDATE "session"
SET
  auth_method = COALESCE(sqlc.narg('auth_method'), auth_method),
  auth_provider = COALESCE(sqlc.narg('auth_provider'), auth_provider),
  expires_at = COALESCE(sqlc.narg('expires_at'), expires_at),
  ip_address = COALESCE(sqlc.narg('ip_address'), ip_address),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id),
  token = COALESCE(sqlc.narg('token'), token),
  user_agent = COALESCE(sqlc.narg('user_agent'), user_agent),
  user_id = COALESCE(sqlc.narg('user_id'), user_id)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteSession :exec
DELETE FROM "session"
WHERE
  id = sqlc.arg('id');