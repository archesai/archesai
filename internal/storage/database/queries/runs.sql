-- name: CreateRun :one
INSERT INTO run (
    organization_id,
    pipeline_id,
    tool_id,
    status,
    progress
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetRun :one
SELECT * FROM run
WHERE id = $1 LIMIT 1;

-- name: ListRuns :many
SELECT * FROM run
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListRunsByOrganization :many
SELECT * FROM run
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRunsByPipeline :many
SELECT * FROM run
WHERE pipeline_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRunsByTool :many
SELECT * FROM run
WHERE tool_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateRun :one
UPDATE run
SET 
    pipeline_id = COALESCE(sqlc.narg(pipeline_id), pipeline_id),
    tool_id = COALESCE(sqlc.narg(tool_id), tool_id),
    status = COALESCE(sqlc.narg(status), status),
    progress = COALESCE(sqlc.narg(progress), progress),
    error = COALESCE(sqlc.narg(error), error),
    started_at = COALESCE(sqlc.narg(started_at), started_at),
    completed_at = COALESCE(sqlc.narg(completed_at), completed_at),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteRun :exec
DELETE FROM run
WHERE id = $1;

-- name: DeleteRunsByPipeline :exec
DELETE FROM run
WHERE pipeline_id = $1;