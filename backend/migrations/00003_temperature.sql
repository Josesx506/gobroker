-- +goose Up
-- +goose StatementBegin
-- Create the temperature readings table
CREATE TABLE IF NOT EXISTS public.temperature_readings (
    id BIGSERIAL PRIMARY KEY,
    location_id TEXT NOT NULL REFERENCES public.device_metadata (location_id) ON DELETE CASCADE,
    device_id TEXT NOT NULL REFERENCES public.device_metadata (device_id) ON DELETE CASCADE,
    value NUMERIC(5, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.temperature_readings;
-- +goose StatementEnd
