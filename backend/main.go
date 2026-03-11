package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Josesx506/gobroker/backend/internal/api"
	"github.com/Josesx506/gobroker/backend/internal/app"
	"github.com/Josesx506/gobroker/backend/internal/broker"
	"github.com/Josesx506/gobroker/backend/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Bacground Broker for new entries
	brkr := broker.NewBroker()
	go brkr.Start()
	go store.PostgresListener(ctx, brkr)

	// REST API logic
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	/** DEV dependencies - unrequired when separate service writes to the DB */
	// Seed the db if data does not exist
	if err := store.SeedDB(app.DB, app.Logger); err != nil {
		app.Logger.Printf("Warning: seed failed: %v", err)
	}
	// Simulate devices sending live data every 2 minutes - 2*time.Minute
	go store.SimulateReadings(ctx, app.DB, 2*time.Minute, app.Logger)

	err = godotenv.Load()
	if err != nil {
		app.Logger.Printf("Error loading .env file")
	}
	port := os.Getenv("PORT")
	defer app.DB.Close()

	router := api.SetupRoutes(app, brkr)
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	app.Logger.Printf("[Main] We are running our api on port %s\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
