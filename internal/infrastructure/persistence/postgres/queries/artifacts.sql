-- name: CreateArtifact :one
INSERT INTO
  artifact (
    id,
    organization_id,
    name,
    description,
    mime_type,
    url,
    credits,
    preview_image,
    producer_id,
    text
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING
  *;

-- name: GetArtifact :one
SELECT
  *
FROM
  artifact
WHERE
  id = $1
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
  $1
OFFSET
  $2;

-- name: ListArtifactsByOrganization :many
SELECT
  *
FROM
  artifact
WHERE
  organization_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: ListArtifactsByProducer :many
SELECT
  *
FROM
  artifact
WHERE
  producer_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;

-- name: UpdateArtifact :one
UPDATE artifact
SET
  name = COALESCE(sqlc.narg (name), name),
  description = COALESCE(sqlc.narg (description), description),
  mime_type = COALESCE(sqlc.narg (mime_type), mime_type),
  url = COALESCE(sqlc.narg (url), url),
  credits = COALESCE(sqlc.narg (credits), credits),
  preview_image = COALESCE(sqlc.narg (preview_image), preview_image),
  text = COALESCE(sqlc.narg (text), text)
WHERE
  id = $1
RETURNING
  *;

-- name: DeleteArtifact :exec
DELETE FROM artifact
WHERE
  id = $1;

-- name: DeleteArtifactsByOrganization :exec
DELETE FROM artifact
WHERE
  organization_id = $1;
