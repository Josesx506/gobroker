-- +goose Up
-- +goose StatementBegin
-- Update the notify function for temperature
CREATE OR REPLACE FUNCTION rtmfuncs.notify_temp_insert()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  -- We broadcast the new reading as a JSON string
  PERFORM pg_notify('temp_events', json_build_object(
    'val', NEW.value,
    'time', NEW.created_at,
    'device', NEW.device_id,
    'location_id', NEW.location_id
  )::text);
  RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS rtmfuncs.notify_temp_insert();
-- +goose StatementEnd
