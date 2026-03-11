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

	err = godotenv.Load()
	if err != nil {
		app.Logger.Printf("Error loading .env file")
	}
	port := os.Getenv("PORT")
	defer app.DB.Close()

	router := api.SetupRoutes(app)
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	app.Logger.Printf("We are running our api on port %s\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
