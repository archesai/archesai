-- name: CreateAccount :one
INSERT INTO
  account (
    id,
    user_id,
    provider_id,
    account_id,
    access_token,
    refresh_token,
    access_token_expires_at,
    refresh_token_expires_at,
    scope,
    id_token,
    password
  )
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    sqlc.narg ('access_token'),
    sqlc.narg ('refresh_token'),
    sqlc.narg ('access_token_expires_at'),
    sqlc.narg ('refresh_token_expires_at'),
    sqlc.narg ('scope'),
    sqlc.narg ('id_token'),
    sqlc.narg ('password')
  )
RETURNING
  *;

-- name: GetAccount :one
SELECT
  *
FROM
  account
WHERE
  id = $1
LIMIT
  1;

-- name: GetAccountByUser :one
SELECT
  *
FROM
  account
WHERE
  user_id = $1
  AND provider_id = $2
LIMIT
  1;

-- name: GetAccountByProviderID :one
SELECT
  *
FROM
  account
WHERE
  provider_id = $1
  AND account_id = $2
LIMIT
  1;

-- name: ListAccounts :many
SELECT
  *
FROM
  account
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: ListAccountsByUser :many
SELECT
  *
FROM
  account
WHERE
  user_id = $1
ORDER BY
  created_at DESC;

-- name: UpdateAccount :one
UPDATE account
SET
  access_token = COALESCE(sqlc.narg (access_token), access_token),
  refresh_token = COALESCE(sqlc.narg (refresh_token), refresh_token),
  access_token_expires_at = COALESCE(
    sqlc.narg (access_token_expires_at),
    access_token_expires_at
  ),
  refresh_token_expires_at = COALESCE(
    sqlc.narg (refresh_token_expires_at),
    refresh_token_expires_at
  ),
  scope = COALESCE(sqlc.narg (scope), scope),
  id_token = COALESCE(sqlc.narg (id_token), id_token),
  password = COALESCE(sqlc.narg (password), password),
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE
  id = $1;

-- name: DeleteAccountsByUser :exec
DELETE FROM account
WHERE
  user_id = $1;
