-- +goose Up
-- +goose StatementBegin
-- Create a channel name constant
CREATE SCHEMA IF NOT EXISTS rtmfuncs;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA rtmfuncs;
-- +goose StatementEnd