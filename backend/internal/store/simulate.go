package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// SimulateReadings inserts one temperature reading per location on every tick
// of interval. It reuses the same day/night temperature curve as the seed data
// so values are consistent. The goroutine exits cleanly when ctx is cancelled.
//
// For production use 10*time.Minute. For testing the LISTEN/NOTIFY pipeline a
// shorter interval (e.g. 5*time.Second) lets you observe live events quickly.
func SimulateReadings(ctx context.Context, db *sql.DB, interval time.Duration, logger *log.Logger) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Printf("[Store] Device simulator started (interval: %s)\n", interval)

	for {
		select {
		case <-ctx.Done():
			logger.Println("[Store] Simulator stopped.")
			return
		case t := <-ticker.C:
			if err := insertReadings(db, t.UTC()); err != nil {
				logger.Printf("[Store] Simulator: insert error: %v\n", err)
			}
		}
	}
}

func insertReadings(db *sql.DB, t time.Time) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO temperature_readings (location_id, device_id, value, created_at)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, loc := range seedLocations {
		value := temperature(t, loc.baseTemp, loc.amplitude)
		if _, err := stmt.Exec(loc.name, "device-"+loc.name, value, t); err != nil {
			return fmt.Errorf("[Store] location %s: %w", loc.name, err)
		}
	}

	return tx.Commit()
}
