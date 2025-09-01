-- name: CreateMember :one
INSERT INTO member (
    user_id,
    organization_id,
    role
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetMember :one
SELECT * FROM member
WHERE id = $1 LIMIT 1;

-- name: GetMemberByUserAndOrg :one
SELECT * FROM member
WHERE user_id = $1 AND organization_id = $2 LIMIT 1;

-- name: ListMembers :many
SELECT * FROM member
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListMembersByOrganization :many
SELECT * FROM member
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListMembersByUser :many
SELECT * FROM member
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateMember :one
UPDATE member
SET 
    role = COALESCE(sqlc.narg(role), role),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteMember :exec
DELETE FROM member
WHERE id = $1;

-- name: DeleteMembersByOrganization :exec
DELETE FROM member
WHERE organization_id = $1;