-- name: CreatePipelineStepDependency :exec
INSERT INTO pipeline_step_to_dependency (
    pipeline_step_id,
    prerequisite_id
) VALUES (
    $1, $2
);

-- name: GetPipelineStepDependencies :many
SELECT 
    psd.pipeline_step_id,
    psd.prerequisite_id,
    ps.pipeline_id
FROM pipeline_step_to_dependency psd
JOIN pipeline_step ps ON ps.id = psd.pipeline_step_id
WHERE ps.pipeline_id = $1;

-- name: GetStepDependencies :many
SELECT prerequisite_id
FROM pipeline_step_to_dependency
WHERE pipeline_step_id = $1;

-- name: GetStepDependents :many
SELECT pipeline_step_id
FROM pipeline_step_to_dependency
WHERE prerequisite_id = $1;

-- name: DeletePipelineStepDependency :exec
DELETE FROM pipeline_step_to_dependency
WHERE pipeline_step_id = $1 AND prerequisite_id = $2;

-- name: DeleteAllStepDependencies :exec
DELETE FROM pipeline_step_to_dependency
WHERE pipeline_step_id = $1 OR prerequisite_id = $1;

-- name: CheckDirectCircularDependency :one
-- Check if adding this dependency would create a direct cycle
SELECT EXISTS (
    SELECT 1 
    FROM pipeline_step_to_dependency 
    WHERE pipeline_step_id = $2 
    AND prerequisite_id = $1
) as would_create_cycle;