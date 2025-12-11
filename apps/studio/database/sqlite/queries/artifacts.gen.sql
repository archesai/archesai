

-- name: CreateArtifact :one
INSERT INTO
  artifact (id, credits, description, mime_type, name, organization_id, preview_image, producer_id, text, url)
VALUES
  (
    $1,
    sqlc.arg('credits'),
    sqlc.narg('description'),
    sqlc.arg('mime_type'),
    sqlc.narg('name'),
    sqlc.arg('organization_id'),
    sqlc.narg('preview_image'),
    sqlc.narg('producer_id'),
    sqlc.narg('text'),
    sqlc.narg('url')
  )
RETURNING
  *;

-- name: GetArtifact :one
SELECT
  *
FROM
  artifact
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListArtifacts :many
SELECT
  *
FROM
  artifact
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountArtifacts :one
SELECT
  COUNT(*)
FROM
  artifact;

-- name: UpdateArtifact :one
UPDATE artifact
SET
  credits = COALESCE(sqlc.narg('credits'), credits),
  description = COALESCE(sqlc.narg('description'), description),
  mime_type = COALESCE(sqlc.narg('mime_type'), mime_type),
  name = COALESCE(sqlc.narg('name'), name),
  organization_id = COALESCE(sqlc.narg('organization_id'), organization_id),
  preview_image = COALESCE(sqlc.narg('preview_image'), preview_image),
  producer_id = COALESCE(sqlc.narg('producer_id'), producer_id),
  text = COALESCE(sqlc.narg('text'), text),
  url = COALESCE(sqlc.narg('url'), url)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteArtifact :exec
DELETE FROM artifact
WHERE
  id = sqlc.arg('id');

-- name: ListArtifactsByOrganization :many
SELECT
  *
FROM
  artifact
WHERE
  organization_id = sqlc.arg('organization_id')
ORDER BY
  created_at DESC;
-- name: ListArtifactsByProducer :many
SELECT
  *
FROM
  artifact
WHERE
  producer_id = sqlc.arg('producer_id')
ORDER BY
  created_at DESC;