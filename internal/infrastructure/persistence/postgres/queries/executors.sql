-- name: CreateExecutor :one
INSERT INTO
  executor (
    id,
    organization_id,
    name,
    description,
    language,
    execute_code,
    dependencies,
    schema_in,
    schema_out,
    extra_files,
    timeout,
    memory_mb,
    cpu_shares,
    env,
    is_active,
    version
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING
  *;

-- name: GetExecutor :one
SELECT
  *
FROM
  executor
WHERE
  id = $1
LIMIT
  1;

-- name: ListExecutors :many
SELECT
  *
FROM
  executor
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: ListExecutorsByOrganization :many
SELECT
  *
FROM
  executor
WHERE
  organization_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateExecutor :one
UPDATE executor
SET
  name = COALESCE(sqlc.narg (name), name),
  description = COALESCE(sqlc.narg (description), description),
  language = COALESCE(sqlc.narg (language), language),
  execute_code = COALESCE(sqlc.narg (execute_code), execute_code),
  dependencies = COALESCE(sqlc.narg (dependencies), dependencies),
  schema_in = COALESCE(sqlc.narg (schema_in), schema_in),
  schema_out = COALESCE(sqlc.narg (schema_out), schema_out),
  extra_files = COALESCE(sqlc.narg (extra_files), extra_files),
  timeout = COALESCE(sqlc.narg (timeout), timeout),
  memory_mb = COALESCE(sqlc.narg (memory_mb), memory_mb),
  cpu_shares = COALESCE(sqlc.narg (cpu_shares), cpu_shares),
  env = COALESCE(sqlc.narg (env), env),
  is_active = COALESCE(sqlc.narg (is_active), is_active),
  version = version + 1,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = $1
RETURNING
  *;

-- name: DeleteExecutor :exec
DELETE FROM executor
WHERE
  id = $1;

-- name: DeleteExecutorsByOrganization :exec
DELETE FROM executor
WHERE
  organization_id = $1;
