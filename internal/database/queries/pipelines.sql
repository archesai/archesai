-- name: CreatePipeline :one
INSERT INTO
  pipeline (id, organization_id, name, description)
VALUES
  ($1, $2, $3, $4)
RETURNING
  *;

-- name: GetPipeline :one
SELECT
  *
FROM
  pipeline
WHERE
  id = $1
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
  $1
OFFSET
  $2;

-- name: ListPipelinesByOrganization :many
SELECT
  *
FROM
  pipeline
WHERE
  organization_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdatePipeline :one
UPDATE pipeline
SET
  name = COALESCE(sqlc.narg (name), name),
  description = COALESCE(sqlc.narg (description), description)
WHERE
  id = $1
RETURNING
  *;

-- name: DeletePipeline :exec
DELETE FROM pipeline
WHERE
  id = $1;

-- name: DeletePipelinesByOrganization :exec
DELETE FROM pipeline
WHERE
  organization_id = $1;
