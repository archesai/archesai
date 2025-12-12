

-- name: CreateOrganization :one
INSERT INTO
  organization (id, billing_email, credits, logo, name, plan, slug, stripe_customer_identifier)
VALUES
  (
    $1,
    sqlc.narg('billing_email'),
    sqlc.arg('credits'),
    sqlc.narg('logo'),
    sqlc.arg('name'),
    sqlc.arg('plan'),
    sqlc.arg('slug'),
    sqlc.arg('stripe_customer_identifier')
  )
RETURNING
  *;

-- name: GetOrganization :one
SELECT
  *
FROM
  organization
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListOrganizations :many
SELECT
  *
FROM
  organization
ORDER BY
  created_at DESC
LIMIT
  sqlc.arg('limit')
OFFSET
  sqlc.arg('offset');

-- name: CountOrganizations :one
SELECT
  COUNT(*)
FROM
  organization;

-- name: UpdateOrganization :one
UPDATE organization
SET
  billing_email = COALESCE(sqlc.narg('billing_email'), billing_email),
  credits = COALESCE(sqlc.narg('credits'), credits),
  logo = COALESCE(sqlc.narg('logo'), logo),
  name = COALESCE(sqlc.narg('name'), name),
  plan = COALESCE(sqlc.narg('plan'), plan),
  slug = COALESCE(sqlc.narg('slug'), slug),
  stripe_customer_identifier = COALESCE(sqlc.narg('stripe_customer_identifier'), stripe_customer_identifier)
WHERE
  id = sqlc.arg('id')
RETURNING
  *;

-- name: DeleteOrganization :exec
DELETE FROM organization
WHERE
  id = sqlc.arg('id');

-- name: GetOrganizationBySlug :one
SELECT
  *
FROM
  organization
WHERE
  slug = sqlc.arg('slug')
LIMIT
  1;
-- name: GetOrganizationByStripeCustomerID :one
SELECT
  *
FROM
  organization
WHERE
  stripe_customer_identifier = sqlc.arg('stripe_customer_identifier')
LIMIT
  1;