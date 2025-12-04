

-- name: CreateInvitation :one
INSERT INTO
  invitation (id, email, expires_at, inviter_id, organization_id, role, status)
VALUES
  (
    $1,
    sqlc.arg('email'),
    sqlc.arg('expires_at'),
    sqlc.arg('inviter_id'),
    sqlc.arg('organization_id'),
    sqlc.arg('role'),
    sqlc.arg('status')
  )
RETURNING
  *;

-- name: GetInvitation :one
SELECT
  *
FROM
  invitation
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListInvitations :many
SELECT
  *
FROM
  invitation
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountInvitations :one
SELECT
  COUNT(*)
FROM
  invitation;

-- name: UpdateInvitation :one
UPDATE invitation
SET
  email = COALESCE(sqlc.narg('email'), email),
  expires_at = COALESCE(sqlc.narg('expires_at'), expires_at),
  inviter_id = COALESCE(sqlc.narg('inviter_id'), inviter_id),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id),
  role = COALESCE(sqlc.narg('role'), role),
  status = COALESCE(sqlc.narg('status'), status)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteInvitation :exec
DELETE FROM invitation
WHERE
  id = sqlc.arg('id');

-- name: ListInvitationsByOrganization :many
SELECT
  *
FROM
  invitation
WHERE
  organization_id = sqlc.arg('organization_id')
ORDER BY
  created_at DESC;
-- name: GetInvitationByEmail :one
SELECT
  *
FROM
  invitation
WHERE
  email = sqlc.arg('email') AND
  organization_id = sqlc.arg('organization_id')
LIMIT
  1;
-- name: ListInvitationsByInviter :many
SELECT
  *
FROM
  invitation
WHERE
  inviter_id = sqlc.arg('inviter_id')
ORDER BY
  created_at DESC;