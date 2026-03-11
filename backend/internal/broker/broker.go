package broker

import "encoding/json"

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
		case s := <-b.closingClients:
			delete(b.clients, s)
		case msg := <-b.notifier:
			// 1. Parse the location_id from the JSON payload
			var data struct {
				LocationID string `json:"location_id"`
			}
			json.Unmarshal([]byte(msg), &data)

			// 2. ONLY send to clients watching this specific location
			for ch, locID := range b.clients {
				if locID == data.LocationID {
					ch <- msg
				}
			}
		}
	}
}

func (b *Broker) Broadcast(msg string) {
	b.notifier <- msg
}
