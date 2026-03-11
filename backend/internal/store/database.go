package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func Open(logger *log.Logger) (*sql.DB, error) {
	// Load .env for local development. In production (Railway, etc.) env vars
	// are injected directly, so skip loading the file if they are already set.
	if os.Getenv("DATABASE_URL") == "" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("DATABASE_URL not set and no .env file found")
		}
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set\n")
	}

	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	// Additional security config
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Second * 60)

	logger.Printf("[Store] Connected to database....\n")
	return db, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)

	// Anonymous function run at the end that wipes global goose
	// state to default behavior
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres") // specify db type
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
