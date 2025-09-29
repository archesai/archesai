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

-- name: GetPipelineStepsWithDependencies :many
SELECT
  ps.id,
  ps.pipeline_id,
  ps.tool_id,
  ps.created_at,
  ps.updated_at,
  COALESCE(
    ARRAY_AGG(DISTINCT psd.prerequisite_id) FILTER (
      WHERE
        psd.prerequisite_id IS NOT NULL
    ),
    ARRAY[]::UUID[]
  ) as dependencies
FROM
  pipeline_step ps
  LEFT JOIN pipeline_step_to_dependency psd ON ps.id = psd.pipeline_step_id
WHERE
  ps.pipeline_id = $1
GROUP BY
  ps.id,
  ps.pipeline_id,
  ps.tool_id,
  ps.created_at,
  ps.updated_at
ORDER BY
  ps.created_at ASC;
