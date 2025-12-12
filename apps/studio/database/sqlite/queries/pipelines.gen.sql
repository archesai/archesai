

-- name: CreatePipeline :one
INSERT INTO
  pipeline (id, description, name, organization_id)
VALUES
  (
    $1,
    sqlc.narg('description'),
    sqlc.narg('name'),
    sqlc.arg('organization_id')
  )
RETURNING
  *;

-- name: GetPipeline :one
SELECT
  *
FROM
  pipeline
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListPipelines :many
SELECT
  *
FROM
  pipeline
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountPipelines :one
SELECT
  COUNT(*)
FROM
  pipeline;

-- name: UpdatePipeline :one
UPDATE pipeline
SET
  description = COALESCE(sqlc.narg('description'), description),
  name = COALESCE(sqlc.narg('name'), name),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeletePipeline :exec
DELETE FROM pipeline
WHERE
  id = sqlc.arg('id');

-- name: ListPipelinesByOrganization :many
SELECT
  *
FROM
  pipeline
WHERE
  organization_id = sqlc.arg('organization_id')
ORDER BY
  created_at DESC;