

-- name: CreatePipelineStep :one
INSERT INTO
  pipeline_step (id, pipeline_id, tool_id)
VALUES
  (
    $1,
    sqlc.arg('pipeline_id'),
    sqlc.arg('tool_id')
  )
RETURNING
  *;

-- name: GetPipelineStep :one
SELECT
  *
FROM
  pipeline_step
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListPipelineSteps :many
SELECT
  *
FROM
  pipeline_step
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountPipelineSteps :one
SELECT
  COUNT(*)
FROM
  pipeline_step;

-- name: UpdatePipelineStep :one
UPDATE pipeline_step
SET
  pipeline_id = COALESCE(sqlc.narg('pipeline_id'), pipeline_id),
  tool_id = COALESCE(sqlc.narg('tool_id'), tool_id)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeletePipelineStep :exec
DELETE FROM pipeline_step
WHERE
  id = sqlc.arg('id');