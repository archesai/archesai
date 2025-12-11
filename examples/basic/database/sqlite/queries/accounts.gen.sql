

-- name: CreateAccount :one
INSERT INTO
  account (id, access_token, access_token_expires_at, account_identifier, id_token, provider, refresh_token, refresh_token_expires_at, scope, user_id)
VALUES
  (
    $1,
    sqlc.narg('access_token'),
    sqlc.narg('access_token_expires_at'),
    sqlc.arg('account_identifier'),
    sqlc.narg('id_token'),
    sqlc.arg('provider'),
    sqlc.narg('refresh_token'),
    sqlc.narg('refresh_token_expires_at'),
    sqlc.narg('scope'),
    sqlc.arg('user_id')
  )
RETURNING
  *;

-- name: GetAccount :one
SELECT
  *
FROM
  account
WHERE
  id = sqlc.arg('id')
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
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountAccounts :one
SELECT
  COUNT(*)
FROM
  account;

-- name: UpdateAccount :one
UPDATE account
SET
  access_token = COALESCE(sqlc.narg('access_token'), access_token),
  access_token_expires_at = COALESCE(sqlc.narg('access_token_expires_at'), access_token_expires_at),
  account_identifier = COALESCE(sqlc.narg('account_identifier'), account_identifier),
  id_token = COALESCE(sqlc.narg('id_token'), id_token),
  provider = COALESCE(sqlc.narg('provider'), provider),
  refresh_token = COALESCE(sqlc.narg('refresh_token'), refresh_token),
  refresh_token_expires_at = COALESCE(sqlc.narg('refresh_token_expires_at'), refresh_token_expires_at),
  scope = COALESCE(sqlc.narg('scope'), scope),
  user_id = COALESCE(sqlc.narg('user_id'), user_id)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE
  id = sqlc.arg('id');

-- name: GetAccountByProvider :one
SELECT
  *
FROM
  account
WHERE
  provider = sqlc.arg('provider') AND
  account_identifier = sqlc.arg('account_identifier')
LIMIT
  1;
-- name: ListAccountsByUserID :many
SELECT
  *
FROM
  account
WHERE
  user_id = sqlc.arg('user_id')
ORDER BY
  created_at DESC;