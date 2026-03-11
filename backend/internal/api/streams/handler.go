package streams

import (
	"fmt"
	"net/http"

	"github.com/Josesx506/gobroker/backend/internal/broker"
)

type StreamHandler struct {
	broker *broker.Broker
}

func NewStreamHandler(broker *broker.Broker) *StreamHandler {
	return &StreamHandler{
		broker: broker,
	}
}

func (sh StreamHandler) HandleLocationTemperatures(w http.ResponseWriter, r *http.Request) {
	locationID := r.URL.Query().Get("location_id")
	if locationID == "" {
		http.Error(w, "location_id is required", http.StatusBadRequest)
		return
	}
	// Set appropriate text header for event stream
	w.Header().Set("Content-Type", "text/event-stream")

	// Create a unique buffered channel for this specific browser connection
	messageChan := make(chan string, 10)

	// Register the client with the broker and defer closing connection
	sh.broker.RegisterClient(locationID, messageChan)
	defer sh.broker.DeregisterClient(messageChan)

	// Listen for data and flush it to the response
	for {
		select {
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}
