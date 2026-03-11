package api

import (
	"net/http"
	"os"

	"github.com/Josesx506/gobroker/backend/internal/api/locations"
	"github.com/Josesx506/gobroker/backend/internal/api/streams"
	"github.com/Josesx506/gobroker/backend/internal/api/temperature"
	"github.com/Josesx506/gobroker/backend/internal/app"
	"github.com/Josesx506/gobroker/backend/internal/broker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func cors(next http.Handler) http.Handler {
	clientURL := os.Getenv("CLIENT_URL")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", clientURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes(app *app.Application, broker *broker.Broker) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors)
	router.Use(middleware.Logger)

	router.Get("/", HandlerRoot)
	router.Get("/health", HandlerHealthChecker)

	// Device metadata
	router.Mount("/locations", locations.LocationsRouter(app))
	// Real-time streaming of temperature data from postgres
	router.Mount("/streams", streams.StreamsRouter(app, broker))
	// API aggregation for historical analytics
	router.Mount("/temperature", temperature.TemperatureRouter(app))

	return router
}
