package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/Josesx506/gobroker/backend/internal/store"
	"github.com/Josesx506/gobroker/backend/migrations"
)

type Application struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewApplication() (*Application, error) {
	// Logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Open db connection
	pgDB, err := store.Open(logger)
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	app := &Application{
		Logger: logger,
		DB:     pgDB,
	}

	return app, nil
}
