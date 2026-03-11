package temperature

import (
	"database/sql"
	"time"
)

// DB connector struct
type PostgresTemperatureStore struct {
	db *sql.DB
}

func NewPostgresTemperatureStore(db *sql.DB) *PostgresTemperatureStore {
	return &PostgresTemperatureStore{db: db}
}

type TemperatureStore interface {
	TemperatureByLocation(locationID string, lookback string) (*ChartResponse, error)
}

// locationID - string for location in DB \n
// lookback - "24h","1wk","1mo","3mo"
func (pts *PostgresTemperatureStore) TemperatureByLocation(locationID string, lookback string) (*ChartResponse, error) {
	var cutoff time.Time
	var truncateUnit string

	now := time.Now().UTC()

	// 1. Determine Lookback and Aggregation Level
	switch lookback {
	case "1wk": // 168 points (7 days * 24h)
		cutoff = now.AddDate(0, 0, -7)
		truncateUnit = "hour"
	case "1mo": // ~30 points
		cutoff = now.AddDate(0, -1, 0)
		truncateUnit = "day"
	case "3mo": // ~90 points
		cutoff = now.AddDate(0, -3, 0)
		truncateUnit = "day"
	default: // 24h from start of today UTC - high details
		cutoff = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		truncateUnit = "minute"
	}

	// 2. The SQL Query - use date_trunc to group data into the desired unit
	query := `
        SELECT 
            date_trunc($1, created_at) AS bucket, 
            ROUND(AVG(value), 2) as avg_value
        FROM temperature_readings 
        WHERE location_id = $2 AND created_at >= $3
        GROUP BY bucket
        ORDER BY bucket ASC;
    `

	rows, err := pts.db.Query(query, truncateUnit, locationID, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*DataPoint
	for rows.Next() {
		dprw := &DataPoint{}
		if err := rows.Scan(&dprw.Time, &dprw.Value); err != nil {
			return nil, err
		}
		results = append(results, dprw)
	}

	response := &ChartResponse{
		Label: locationID,
		Data:  results,
	}

	return response, nil
}
