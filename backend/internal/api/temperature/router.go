package temperature

import (
	"github.com/Josesx506/gobroker/backend/internal/app"
	"github.com/go-chi/chi/v5"
)

func TemperatureRouter(app *app.Application) chi.Router {
	r := chi.NewRouter()

	store := NewPostgresTemperatureStore(app.DB)
	handler := NewTemperatureHandler(store, app.Logger)

	r.Get("/{location_id}", handler.HandleGetTemperatureByLocation)

	return r
}
