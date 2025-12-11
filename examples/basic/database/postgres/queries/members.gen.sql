

-- name: CreateMember :one
INSERT INTO
  member (id, organization_id, role, user_id)
VALUES
  (
    $1,
    sqlc.arg('organization_id'),
    sqlc.arg('role'),
    sqlc.arg('user_id')
  )
RETURNING
  *;

-- name: GetMember :one
SELECT
  *
FROM
  member
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListMembers :many
SELECT
  *
FROM
  member
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountMembers :one
SELECT
  COUNT(*)
FROM
  member;

-- name: UpdateMember :one
UPDATE member
SET
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id),
  role = COALESCE(sqlc.narg('role'), role),
  user_id = COALESCE(sqlc.narg('user_id'), user_id)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteMember :exec
DELETE FROM member
WHERE
  id = sqlc.arg('id');

-- name: ListMembersByOrganization :many
SELECT
  *
FROM
  member
WHERE
  organization_id = sqlc.arg('organization_id')
ORDER BY
  created_at DESC;
-- name: ListMembersByUser :many
SELECT
  *
FROM
  member
WHERE
  user_id = sqlc.arg('user_id')
ORDER BY
  created_at DESC;
-- name: GetMemberByUserAndOrganization :one
SELECT
  *
FROM
  member
WHERE
  user_id = sqlc.arg('user_id') AND
  organization_id = sqlc.arg('organization_id')
LIMIT
  1;