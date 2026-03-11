package locations

import (
	"database/sql"
	"fmt"
)

type PostgresLocationStore struct {
	db *sql.DB
}

func NewPostgresLocationStore(db *sql.DB) *PostgresLocationStore {
	return &PostgresLocationStore{db: db}
}

type LocationStore interface {
	AllLocations() ([]*DeviceLocation, error)
}

func (s *PostgresLocationStore) AllLocations() ([]*DeviceLocation, error) {
	rows, err := s.db.Query(`
		SELECT id, location_id, device_id, longitude, latitude
		FROM device_metadata
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: all locations: %w", err)
	}
	defer rows.Close()

	var locations []*DeviceLocation
	for rows.Next() {
		loc := &DeviceLocation{}
		if err := rows.Scan(&loc.ID, &loc.LocationID, &loc.DeviceID, &loc.Longitude, &loc.Latitude); err != nil {
			return nil, fmt.Errorf("store: scan location: %w", err)
		}
		locations = append(locations, loc)
	}

	return locations, rows.Err()
}
