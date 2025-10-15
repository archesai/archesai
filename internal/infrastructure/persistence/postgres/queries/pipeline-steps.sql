-- name: CreatePipelineStep :one
INSERT INTO
  pipeline_step (id, pipeline_id, tool_id)
VALUES
  ($1, $2, $3)
RETURNING
  *;

-- name: GetPipelineStep :one
SELECT
  *
FROM
  pipeline_step
WHERE
  id = $1
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
  $1
OFFSET
  $2;

-- name: ListPipelineStepsByPipeline :many
SELECT
  *
FROM
  pipeline_step
WHERE
  pipeline_id = $1
ORDER BY
  created_at ASC
LIMIT
  $2
OFFSET
  $3;

-- name: ListPipelineStepsByTool :many
SELECT
  *
FROM
  pipeline_step
WHERE
  tool_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdatePipelineStep :one
UPDATE pipeline_step
SET
  tool_id = COALESCE(sqlc.narg (tool_id), tool_id)
WHERE
  id = $1
RETURNING
  *;

-- name: DeletePipelineStep :exec
DELETE FROM pipeline_step
WHERE
  id = $1;

-- name: DeletePipelineStepsByPipeline :exec
DELETE FROM pipeline_step
WHERE
  pipeline_id = $1;

-- name: CountPipelineSteps :one
SELECT
  COUNT(*) as count
FROM
  pipeline_step
WHERE
  pipeline_id = $1;
