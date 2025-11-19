-- name: CreateTool :one
INSERT INTO
  tool (
    id,
    organization_id,
    name,
    description,
    input_mime_type,
    output_mime_type
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  *;

-- name: GetTool :one
SELECT
  *
FROM
  tool
WHERE
  id = $1
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
  $1
OFFSET
  $2;

-- name: ListToolsByOrganization :many
SELECT
  *
FROM
  tool
WHERE
  organization_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateTool :one
UPDATE tool
SET
  name = COALESCE(sqlc.narg (name), name),
  description = COALESCE(sqlc.narg (description), description),
  input_mime_type = COALESCE(sqlc.narg (input_mime_type), input_mime_type),
  output_mime_type = COALESCE(sqlc.narg (output_mime_type), output_mime_type)
WHERE
  id = $1
RETURNING
  *;

-- name: DeleteTool :exec
DELETE FROM tool
WHERE
  id = $1;

-- name: DeleteToolsByOrganization :exec
DELETE FROM tool
WHERE
  organization_id = $1;
