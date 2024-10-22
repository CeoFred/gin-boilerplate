package database

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*
var migrationFiles embed.FS

func RunManualMigration(dbURL string) {
	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Use the embedded migration files
	migrations, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		log.Fatal("Failed to create sub filesystem for migrations:", err)
	}

	d, err := iofs.New(migrations, ".")
	if err != nil {
		log.Fatal("Failed to create iofs driver:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create postgres driver:", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Completed DB migration")
}
