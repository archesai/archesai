-- name: CreateOrganization :one
INSERT INTO organization (
    name,
    billing_email,
    plan,
    credits,
    logo,
    metadata,
    stripe_customer_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organization
WHERE id = $1 LIMIT 1;

-- name: ListOrganizations :many
SELECT * FROM organization
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrganization :one
UPDATE organization
SET 
    name = COALESCE(sqlc.narg(name), name),
    billing_email = COALESCE(sqlc.narg(billing_email), billing_email),
    plan = COALESCE(sqlc.narg(plan), plan),
    credits = COALESCE(sqlc.narg(credits), credits),
    logo = COALESCE(sqlc.narg(logo), logo),
    metadata = COALESCE(sqlc.narg(metadata), metadata),
    stripe_customer_id = COALESCE(sqlc.narg(stripe_customer_id), stripe_customer_id)
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organization
WHERE id = $1;