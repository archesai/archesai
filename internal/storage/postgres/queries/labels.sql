-- name: CreateLabel :one
INSERT INTO label (
    organization_id,
    name
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetLabel :one
SELECT * FROM label
WHERE id = $1 LIMIT 1;

-- name: GetLabelByName :one
SELECT * FROM label
WHERE organization_id = $1 AND name = $2 LIMIT 1;

-- name: ListLabels :many
SELECT * FROM label
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListLabelsByOrganization :many
SELECT * FROM label
WHERE organization_id = $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: UpdateLabel :one
UPDATE label
SET 
    name = COALESCE(sqlc.narg(name), name)
WHERE id = $1
RETURNING *;

-- name: DeleteLabel :exec
DELETE FROM label
WHERE id = $1;

-- name: DeleteLabelsByOrganization :exec
DELETE FROM label
WHERE organization_id = $1;