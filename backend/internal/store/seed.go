package store

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

type locationSeed struct {
	name      string
	baseTemp  float64
	amplitude float64
	longitude float64
	latitude  float64
}

var seedLocations = []locationSeed{
	{"alpha", 15.0, 12.5, -112.0740, 33.4484}, // Phoenix
	{"beta", 17.0, 10.0, -110.9747, 32.2226},  // Tucson
	{"charlie", 14.5, 14.0, -111.6513, 35.1983}, // Flagstaff
	{"delta", 16.0, 11.5, -111.7610, 34.8697},  // Sedona
	{"echo", 18.0, 13.0, -111.8315, 33.4152},   // Mesa
}

func SeedDB(db *sql.DB, logger *log.Logger) error {
	if err := seedMetadata(db, logger); err != nil {
		return err
	}
	if err := seedReadings(db, logger); err != nil {
		return err
	}
	return nil
}

func seedMetadata(db *sql.DB, logger *log.Logger) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM device_metadata`).Scan(&count); err != nil {
		return fmt.Errorf("[Store] seed: metadata check failed: %w", err)
	}
	if count > 0 {
		logger.Println("[Store] Device metadata already present, skipping...")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("[Store] seed: begin metadata transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO device_metadata (location_id, device_id, longitude, latitude)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("[Store] seed: prepare metadata statement: %w", err)
	}
	defer stmt.Close()

	for _, loc := range seedLocations {
		if _, err := stmt.Exec(loc.name, "device-"+loc.name, loc.longitude, loc.latitude); err != nil {
			return fmt.Errorf("[Store] seed: insert metadata %s: %w", loc.name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("[Store] seed: commit metadata: %w", err)
	}

	logger.Println("[Store] Device metadata seeded.")
	return nil
}

func seedReadings(db *sql.DB, logger *log.Logger) error {
	var count int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM temperature_readings
		WHERE location_id IN ('alpha','beta','charlie','delta','echo')
	`).Scan(&count); err != nil {
		return fmt.Errorf("[Store] seed: readings check failed: %w", err)
	}
	if count > 0 {
		logger.Println("[Store] Temperature readings already present, skipping...")
		return nil
	}

	logger.Println("[Store] Seeding temperature readings...")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("[Store] seed: begin readings transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO temperature_readings (location_id, device_id, value, created_at)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("[Store] seed: prepare readings statement: %w", err)
	}
	defer stmt.Close()

	interval := 10 * time.Minute
	start := time.Now().UTC().Add(-30 * 24 * time.Hour)
	end := time.Now().UTC()

	for _, loc := range seedLocations {
		for ts := start; ts.Before(end); ts = ts.Add(interval) {
			value := temperature(ts, loc.baseTemp, loc.amplitude)
			if _, err := stmt.Exec(loc.name, "device-"+loc.name, value, ts); err != nil {
				return fmt.Errorf("[Store] seed: insert reading %s at %v: %w", loc.name, ts, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("[Store] seed: commit readings: %w", err)
	}

	logger.Println("[Store] Temperature readings seeded.")
	return nil
}

// temperature returns a value that rises smoothly from baseTemp during the
// night to baseTemp+amplitude at solar noon (12:00 UTC), using a sine curve
// between 06:00 and 18:00 UTC. A small random jitter (+/- 0.5°C) is added
// to avoid perfectly flat nighttime readings.
func temperature(t time.Time, baseTemp, amplitude float64) float64 {
	hour := float64(t.Hour()) + float64(t.Minute())/60.0
	var delta float64
	if hour >= 6 && hour <= 18 {
		delta = amplitude * math.Sin(math.Pi*(hour-6)/12)
	}
	jitter := (rand.Float64() - 0.5)
	value := baseTemp + delta + jitter
	return math.Round(value*100) / 100
}
