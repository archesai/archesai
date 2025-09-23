-- +goose Up
-- +goose StatementBegin
-- Drop metadata column from organization table
ALTER TABLE organization DROP COLUMN IF EXISTS metadata;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Re-add metadata column if rolling back
ALTER TABLE organization ADD COLUMN metadata TEXT;
-- +goose StatementEnd
