-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.device_metadata (
    id          SERIAL PRIMARY KEY,
    location_id TEXT NOT NULL UNIQUE,
    device_id   TEXT NOT NULL UNIQUE,
    longitude   NUMERIC(9, 6) NOT NULL,
    latitude    NUMERIC(8, 6) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.device_metadata;
-- +goose StatementEnd
