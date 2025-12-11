

-- name: CreateRun :one
INSERT INTO
  run (id, completed_at, error, organization_id, pipeline_id, progress, started_at, status, tool_id)
VALUES
  (
    $1,
    sqlc.narg('completed_at'),
    sqlc.narg('error'),
    sqlc.arg('organization_id'),
    sqlc.arg('pipeline_id'),
    sqlc.arg('progress'),
    sqlc.narg('started_at'),
    sqlc.arg('status'),
    sqlc.arg('tool_id')
  )
RETURNING
  *;

-- name: GetRun :one
SELECT
  *
FROM
  run
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListRuns :many
SELECT
  *
FROM
  run
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountRuns :one
SELECT
  COUNT(*)
FROM
  run;

-- name: UpdateRun :one
UPDATE run
SET
  completed_at = COALESCE(sqlc.narg('completed_at'), completed_at),
  error = COALESCE(sqlc.narg('error'), error),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id),
  pipeline_id = COALESCE(sqlc.narg('pipeline_id'), pipeline_id),
  progress = COALESCE(sqlc.narg('progress'), progress),
  started_at = COALESCE(sqlc.narg('started_at'), started_at),
  status = COALESCE(sqlc.narg('status'), status),
  tool_id = COALESCE(sqlc.narg('tool_id'), tool_id)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteRun :exec
DELETE FROM run
WHERE
  id = sqlc.arg('id');

-- name: ListRunsByPipeline :many
SELECT
  *
FROM
  run
WHERE
  pipeline_id = sqlc.arg('pipeline_id')
ORDER BY
  created_at DESC;
-- name: ListRunsByOrganization :many
SELECT
  *
FROM
  run
WHERE
  organization_id = sqlc.arg('organization_id')
ORDER BY
  created_at DESC;
-- name: ListRunsByTool :many
SELECT
  *
FROM
  run
WHERE
  tool_id = sqlc.arg('tool_id')
ORDER BY
  created_at DESC;