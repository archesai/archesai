

-- name: CreateTodo :one
INSERT INTO
  todo (id, completed, title)
VALUES
  (
    $1,
    sqlc.arg('completed'),
    sqlc.arg('title')
  )
RETURNING
  *;

-- name: GetTodo :one
SELECT
  *
FROM
  todo
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListTodos :many
SELECT
  *
FROM
  todo
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountTodos :one
SELECT
  COUNT(*)
FROM
  todo;

-- name: UpdateTodo :one
UPDATE todo
SET
  completed = COALESCE(sqlc.narg('completed'), completed),
  title = COALESCE(sqlc.narg('title'), title)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteTodo :exec
DELETE FROM todo
WHERE
  id = sqlc.arg('id');