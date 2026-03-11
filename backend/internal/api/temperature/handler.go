package temperature

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type TemperatureHandler struct {
	store  TemperatureStore
	logger *log.Logger
}

func NewTemperatureHandler(store TemperatureStore, logger *log.Logger) *TemperatureHandler {
	return &TemperatureHandler{
		store:  store,
		logger: logger,
	}
}

func (th *TemperatureHandler) HandleGetTemperatureByLocation(w http.ResponseWriter, r *http.Request) {
	locationID := chi.URLParam(r, "location_id")
	lookback := r.URL.Query().Get("lookback")
	if locationID == "" {
		http.Error(w, "Invalid Location", http.StatusBadRequest)
		return
	}
	if lookback == "" {
		lookback = "24h"
	}

	readings, err := th.store.TemperatureByLocation(locationID, lookback)
	if err != nil {
		th.logger.Printf("[Temperature] Read DB error: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Note: You can't call http.Error here if the header was already sent
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(readings); err != nil {
		th.logger.Printf("[Temperature] JSON encoding error: %v\n", err)
		return
	}
}
