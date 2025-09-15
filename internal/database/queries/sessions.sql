-- name: CreateSession :one
INSERT INTO
  session (
    id,
    user_id,
    token,
    expires_at,
    active_organization_id,
    ip_address,
    user_agent
  )
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    sqlc.narg ('active_organization_id'),
    sqlc.narg ('ip_address'),
    sqlc.narg ('user_agent')
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
  active_organization_id = COALESCE(
    sqlc.narg (active_organization_id),
    active_organization_id
  )
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
