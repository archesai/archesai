-- name: CreateApiToken :one
INSERT INTO api_token (
    id,
    user_id,
    name,
    key,
    prefix,
    enabled,
    expires_at,
    permissions,
    rate_limit_enabled,
    rate_limit_max,
    rate_limit_time_window,
    refill_amount,
    refill_interval,
    remaining,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: GetApiToken :one
SELECT * FROM api_token
WHERE id = $1 LIMIT 1;

-- name: GetApiTokenByKey :one
SELECT * FROM api_token
WHERE key = $1 LIMIT 1;

-- name: ListApiTokens :many
SELECT * FROM api_token
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListApiTokensByUser :many
SELECT * FROM api_token
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateApiToken :one
UPDATE api_token
SET 
    name = COALESCE(sqlc.narg(name), name),
    enabled = COALESCE(sqlc.narg(enabled), enabled),
    expires_at = COALESCE(sqlc.narg(expires_at), expires_at),
    permissions = COALESCE(sqlc.narg(permissions), permissions),
    rate_limit_enabled = COALESCE(sqlc.narg(rate_limit_enabled), rate_limit_enabled),
    rate_limit_max = COALESCE(sqlc.narg(rate_limit_max), rate_limit_max),
    rate_limit_time_window = COALESCE(sqlc.narg(rate_limit_time_window), rate_limit_time_window),
    metadata = COALESCE(sqlc.narg(metadata), metadata)
WHERE id = $1
RETURNING *;

-- name: DeleteApiToken :exec
DELETE FROM api_token
WHERE id = $1;

-- name: DeleteApiTokensByUser :exec
DELETE FROM api_token
WHERE user_id = $1;