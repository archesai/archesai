

-- name: CreateTool :one
INSERT INTO
  tool (id, description, input_mime_type, name, organization_id, output_mime_type)
VALUES
  (
    $1,
    sqlc.arg('description'),
    sqlc.arg('input_mime_type'),
    sqlc.arg('name'),
    sqlc.arg('organization_id'),
    sqlc.arg('output_mime_type')
  )
RETURNING
  *;

-- name: GetTool :one
SELECT
  *
FROM
  tool
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListTools :many
SELECT
  *
FROM
  tool
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountTools :one
SELECT
  COUNT(*)
FROM
  tool;

-- name: UpdateTool :one
UPDATE tool
SET
  description = COALESCE(sqlc.narg('description'), description),
  input_mime_type = COALESCE(sqlc.narg('input_mime_type'), input_mime_type),
  name = COALESCE(sqlc.narg('name'), name),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id),
  output_mime_type = COALESCE(sqlc.narg('output_mime_type'), output_mime_type)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteTool :exec
DELETE FROM tool
WHERE
  id = sqlc.arg('id');

-- name: ListToolsByOrganization :many
SELECT
  *
FROM
  tool
WHERE
  organization_id = sqlc.arg('organization_id')
ORDER BY
  created_at DESC;