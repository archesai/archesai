-- name: CreateInvitation :one
INSERT INTO invitation (
    id,
    organization_id,
    inviter_id,
    email,
    role,
    expires_at,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetInvitation :one
SELECT * FROM invitation
WHERE id = $1 LIMIT 1;

-- name: GetInvitationByEmail :one
SELECT * FROM invitation
WHERE organization_id = $1 AND email = $2 LIMIT 1;

-- name: ListInvitations :many
SELECT * FROM invitation
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListInvitationsByOrganization :many
SELECT * FROM invitation
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListInvitationsByInviter :many
SELECT * FROM invitation
WHERE inviter_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateInvitation :one
UPDATE invitation
SET 
    email = COALESCE(sqlc.narg(email), email),
    role = COALESCE(sqlc.narg(role), role),
    expires_at = COALESCE(sqlc.narg(expires_at), expires_at),
    status = COALESCE(sqlc.narg(status), status),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteInvitation :exec
DELETE FROM invitation
WHERE id = $1;

-- name: DeleteInvitationsByOrganization :exec
DELETE FROM invitation
WHERE organization_id = $1;