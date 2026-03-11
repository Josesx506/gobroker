package streams

import (
	"github.com/Josesx506/gobroker/backend/internal/app"
	"github.com/Josesx506/gobroker/backend/internal/broker"
	"github.com/go-chi/chi/v5"
)

func StreamsRouter(app *app.Application, broker *broker.Broker) chi.Router {
	r := chi.NewRouter()

	handler := NewStreamHandler(broker)

	r.Get("/", handler.HandleLocationTemperatures)

	return r
}
