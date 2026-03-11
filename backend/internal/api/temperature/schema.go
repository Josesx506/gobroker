package temperature

import "time"

type TempReadings struct {
	ID         int       `json:"id,omitempty"`
	LocationID string    `json:"location_id,omitempty"`
	DeviceID   string    `json:"device_id,omitempty"`
	Value      float64   `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
}

type DataPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"val"`
}

type ChartResponse struct {
	Label string       `json:"label"`
	Data  []*DataPoint `json:"data"`
}
