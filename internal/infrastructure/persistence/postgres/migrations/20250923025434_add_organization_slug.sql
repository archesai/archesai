-- +goose Up
-- +goose StatementBegin
-- Add slug column to organization table
ALTER TABLE organization
ADD COLUMN slug VARCHAR(50);

-- Create an index on slug for faster lookups
CREATE UNIQUE INDEX idx_organization_slug ON organization(slug);

-- Update existing organizations with a generated slug (if any exist)
UPDATE organization
SET slug = LOWER(REGEXP_REPLACE(REGEXP_REPLACE(name, '[^a-zA-Z0-9]+', '-', 'g'), '^-+|-+$', '', 'g'))
WHERE slug IS NULL;

-- Make slug NOT NULL after populating it
ALTER TABLE organization
ALTER COLUMN slug SET NOT NULL;

-- Add check constraint to ensure slug format
ALTER TABLE organization
ADD CONSTRAINT organization_slug_format CHECK (slug ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the constraint first
ALTER TABLE organization DROP CONSTRAINT IF EXISTS organization_slug_format;

-- Drop the index
DROP INDEX IF EXISTS idx_organization_slug;

-- Drop the column
ALTER TABLE organization DROP COLUMN IF EXISTS slug;
-- +goose StatementEnd
