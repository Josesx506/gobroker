-- +goose Up
-- +goose StatementBegin
-- Index on created_at makes the TTL DELETE a fast range scan rather than a
-- full table scan, keeping the per-insert overhead negligible.
CREATE INDEX IF NOT EXISTS idx_temperature_readings_created_at
    ON public.temperature_readings (created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION rtmfuncs.ttl_temperature_readings()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM public.temperature_readings
    WHERE created_at < NOW() - INTERVAL '3 months';
    RETURN NULL;
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- FOR EACH STATEMENT fires once per INSERT call regardless of row count,
-- so bulk inserts (e.g. seeding) only trigger one cleanup pass.
CREATE OR REPLACE TRIGGER temperature_readings_ttl
    AFTER INSERT ON public.temperature_readings
    FOR EACH STATEMENT
    EXECUTE FUNCTION rtmfuncs.ttl_temperature_readings();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS temperature_readings_ttl ON public.temperature_readings;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS rtmfuncs.ttl_temperature_readings();
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_temperature_readings_created_at;
-- +goose StatementEnd
