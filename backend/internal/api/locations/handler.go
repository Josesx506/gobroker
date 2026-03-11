package locations

import (
	"encoding/json"
	"log"
	"net/http"
)

type LocationHandler struct {
	store  LocationStore
	logger *log.Logger
}

func NewLocationHandler(store LocationStore, logger *log.Logger) *LocationHandler {
	return &LocationHandler{store: store, logger: logger}
}

func (lh *LocationHandler) HandleGetLocations(w http.ResponseWriter, r *http.Request) {
	locs, err := lh.store.AllLocations()
	if err != nil {
		lh.logger.Printf("[Locations] DB error: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locs); err != nil {
		lh.logger.Printf("[Locations] JSON encoding error: %v\n", err)
	}
}
