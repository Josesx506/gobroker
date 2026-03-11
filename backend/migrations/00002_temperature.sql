-- +goose Up
-- +goose StatementBegin
-- Create the temperature readings table
CREATE TABLE IF NOT EXISTS public.temperature_readings (
    id SERIAL PRIMARY KEY,
    location_id TEXT NOT NULL,
    device_id TEXT NOT NULL,
    value NUMERIC(5, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.temperature_readings;
-- +goose StatementEnd
