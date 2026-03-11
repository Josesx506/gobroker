package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Josesx506/gobroker/backend/internal/broker"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func PostgresListener(ctx context.Context, broker *broker.Broker) {
	/**
	Postgres event listener with exponential backoff
	*/
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Printf("DATABASE_URL environment variable is not set\n")
	}

	minWait := 1 * time.Second
	maxWait := 60 * time.Second
	currentWait := minWait

	for {
		err := runListener(ctx, dbURL, broker)
		if err != nil {
			log.Printf("Listener failed: %v. Retrying in %v...", err, currentWait)

			select {
			case <-time.After(currentWait):
				// Exponential increase
				currentWait *= 2
				if currentWait > maxWait {
					currentWait = maxWait
				}
			case <-ctx.Done():
				return
			}
			continue
		}

		// If runListener returns nil, it means it closed gracefully
		// Reset the wait time for the next potential failure
		currentWait = minWait
	}
}

func runListener(ctx context.Context, dbURL string, brokerChn *broker.Broker) error {
	// 1. Establish connection
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	defer conn.Close(ctx)

	// 2. Start Listening
	_, err = conn.Exec(ctx, "LISTEN temp_events")
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}

	log.Println("Successfully connected. Listening for Postgres notifications...")

	for {
		// 3. Block and wait for data
		notification, err := conn.WaitForNotification(ctx)
		if err != nil { // connection was lost
			return fmt.Errorf("notification error: %w", err)
		}

		// 4. Send to broker
		brokerChn.Broadcast(string(notification.Payload))
	}
}
