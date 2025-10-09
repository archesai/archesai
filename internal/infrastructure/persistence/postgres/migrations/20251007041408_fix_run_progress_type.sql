-- +goose Up
-- +goose StatementBegin
-- Change progress column from DOUBLE PRECISION to INTEGER to match entity definition
-- Progress represents a percentage (0-100), so integer is more appropriate
ALTER TABLE run ALTER COLUMN progress TYPE INTEGER USING progress::integer;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert progress column back to DOUBLE PRECISION
ALTER TABLE run ALTER COLUMN progress TYPE DOUBLE PRECISION;
-- +goose StatementEnd
