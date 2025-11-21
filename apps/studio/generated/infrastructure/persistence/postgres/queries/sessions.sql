-- name: CreateSession :one
INSERT INTO
  session (
    id,
    user_id,
    token,
    expires_at,
    organization_id,
    ip_address,
    user_agent,
    auth_method,
    auth_provider,
    created_at,
    updated_at
  )
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    sqlc.narg ('organization_id'),
    sqlc.narg ('ip_address'),
    sqlc.narg ('user_agent'),
    sqlc.narg ('auth_method'),
    sqlc.narg ('auth_provider'),
    NOW(),
    NOW()
  )
RETURNING
  *;

-- name: GetSession :one
SELECT
  *
FROM
  session
WHERE
  id = $1
LIMIT
  1;

-- name: GetSessionByToken :one
SELECT
  *
FROM
  session
WHERE
  token = $1
LIMIT
  1;

-- name: ListSessions :many
SELECT
  *
FROM
  session
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: ListSessionsByUser :many
SELECT
  *
FROM
  session
WHERE
  user_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateSession :one
UPDATE session
SET
  expires_at = COALESCE(sqlc.narg (expires_at), expires_at),
  organization_id = COALESCE(
    sqlc.narg (organization_id),
    organization_id
  ),
  auth_method = COALESCE(sqlc.narg (auth_method), auth_method),
  auth_provider = COALESCE(sqlc.narg (auth_provider), auth_provider),
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  *;

-- name: DeleteSession :exec
DELETE FROM session
WHERE
  id = $1;

-- name: DeleteSessionsByUser :exec
DELETE FROM session
WHERE
  user_id = $1;
