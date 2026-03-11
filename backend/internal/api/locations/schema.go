package locations

type DeviceLocation struct {
	ID         int     `json:"id"`
	LocationID string  `json:"location_id"`
	DeviceID   string  `json:"device_id"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
}
