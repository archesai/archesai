

-- name: CreateExecutor :one
INSERT INTO
  executor (id, cpu_shares, dependencies, description, env, execute_code, extra_files, is_active, language, memory_mb, name, organization_id, schema_in, schema_out, timeout, version)
VALUES
  (
    $1,
    sqlc.arg('cpu_shares'),
    sqlc.narg('dependencies'),
    sqlc.arg('description'),
    sqlc.narg('env'),
    sqlc.arg('execute_code'),
    sqlc.narg('extra_files'),
    sqlc.arg('is_active'),
    sqlc.arg('language'),
    sqlc.arg('memory_mb'),
    sqlc.arg('name'),
    sqlc.arg('organization_id'),
    sqlc.narg('schema_in'),
    sqlc.narg('schema_out'),
    sqlc.arg('timeout'),
    sqlc.arg('version')
  )
RETURNING
  *;

-- name: GetExecutor :one
SELECT
  *
FROM
  executor
WHERE
  id = sqlc.arg('id')
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
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountExecutors :one
SELECT
  COUNT(*)
FROM
  executor;

-- name: UpdateExecutor :one
UPDATE executor
SET
  cpu_shares = COALESCE(sqlc.narg('cpu_shares'), cpu_shares),
  dependencies = COALESCE(sqlc.narg('dependencies'), dependencies),
  description = COALESCE(sqlc.narg('description'), description),
  env = COALESCE(sqlc.narg('env'), env),
  execute_code = COALESCE(sqlc.narg('execute_code'), execute_code),
  extra_files = COALESCE(sqlc.narg('extra_files'), extra_files),
  is_active = COALESCE(sqlc.narg('is_active'), is_active),
  language = COALESCE(sqlc.narg('language'), language),
  memory_mb = COALESCE(sqlc.narg('memory_mb'), memory_mb),
  name = COALESCE(sqlc.narg('name'), name),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id),
  schema_in = COALESCE(sqlc.narg('schema_in'), schema_in),
  schema_out = COALESCE(sqlc.narg('schema_out'), schema_out),
  timeout = COALESCE(sqlc.narg('timeout'), timeout),
  version = COALESCE(sqlc.narg('version'), version)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteExecutor :exec
DELETE FROM executor
WHERE
  id = sqlc.arg('id');

-- name: ListExecutorsByOrganization :many
SELECT
  *
FROM
  executor
WHERE
  organization_id = sqlc.arg('organization_id')
ORDER BY
  created_at DESC;