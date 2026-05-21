-- +goose Up
-- +goose StatementBegin
-- Attach the trigger
CREATE OR REPLACE TRIGGER temp_notify_insert
AFTER INSERT ON public.temperature_readings
FOR EACH ROW 
EXECUTE FUNCTION rtmfuncs.notify_temp_insert();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS temp_notify_insert ON public.temperature_readings;
-- +goose StatementEnd