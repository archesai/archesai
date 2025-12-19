

-- name: CreateLabel :one
INSERT INTO
  label (id, name, organization_id)
VALUES
  (
    $1,
    sqlc.arg('name'),
    sqlc.arg('organization_id')
  )
RETURNING
  *;

-- name: GetLabel :one
SELECT
  *
FROM
  label
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListLabels :many
SELECT
  *
FROM
  label
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountLabels :one
SELECT
  COUNT(*)
FROM
  label;

-- name: UpdateLabel :one
UPDATE label
SET
  name = COALESCE(sqlc.narg('name'), name),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteLabel :exec
DELETE FROM label
WHERE
  id = sqlc.arg('id');

-- name: ListLabelsByOrganization :many
SELECT
  *
FROM
  label
WHERE
  organization_id = sqlc.arg('organization_id')
ORDER BY
  created_at DESC;
-- name: GetLabelByName :one
SELECT
  *
FROM
  label
WHERE
  name = sqlc.arg('name') AND
  organization_id = sqlc.arg('organization_id')
LIMIT
  1;