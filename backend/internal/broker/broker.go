package broker

import (
	"encoding/json"
	"log"
)

type Client struct {
	locationID string
	send       chan string // Each client gets their own channel
}

type Broker struct {
	clients        map[chan string]string // Map of channel to locationID
	newClients     chan Client
	closingClients chan chan string
	notifier       chan string
}

func NewBroker() *Broker {
	broker := &Broker{
		clients:        make(map[chan string]string),
		newClients:     make(chan Client),
		closingClients: make(chan chan string),
		notifier:       make(chan string),
	}

	return broker
}

// Broker methods
func (b *Broker) Start() {
	for {
		select {
		case client := <-b.newClients:
			b.clients[client.send] = client.locationID
			log.Printf("[Broker] New client registered for location: %s (Total clients: %d)",
				client.locationID, len(b.clients))

		// Explicitly close the channel to free resources
		case s := <-b.closingClients:
			locID := b.clients[s]
			delete(b.clients, s)
			close(s)
			log.Printf("[Broker] Client disconnected from location: %s (Total clients: %d)",
				locID, len(b.clients))

		case msg := <-b.notifier:
			// 1. Decode just enough to get the location_id
			var data struct {
				LocationID string `json:"location_id"`
			}
			// Log and skip malformed JSON from Postgres
			if err := json.Unmarshal([]byte(msg), &data); err != nil {
				log.Printf("[Broker] Error unmarshaling Postgres notification: %v", err)
				continue
			}

			// 2. Fan-out to matching clients
			for ch, locID := range b.clients {
				if locID == data.LocationID {
					// 3. Non-blocking send.  If the client's buffer is full, we skip them to keep the broker fast
					select {
					case ch <- msg:
					default: // Optional: log the lagging client
						log.Printf("[Broker] Warning: Dropping message for slow client at %s", locID)

					}
				}
			}
		}
	}
}

func (b *Broker) Broadcast(msg string) {
	b.notifier <- msg
}

func (b *Broker) RegisterClient(locationID string, msgChannel chan string) {
	b.newClients <- Client{
		locationID: locationID,
		send:       msgChannel,
	}
}

func (b *Broker) DeregisterClient(msgChannel chan string) {
	b.closingClients <- msgChannel
}
