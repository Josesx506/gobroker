package locations

import (
	"github.com/Josesx506/gobroker/backend/internal/app"
	"github.com/go-chi/chi/v5"
)

func LocationsRouter(app *app.Application) chi.Router {
	r := chi.NewRouter()

	store := NewPostgresLocationStore(app.DB)
	handler := NewLocationHandler(store, app.Logger)

	r.Get("/", handler.HandleGetLocations)

	return r
}
