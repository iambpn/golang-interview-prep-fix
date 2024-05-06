package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	_ "github.com/joho/godotenv/autoload"

	"log"
	"net/http"

	"github.com/matthewjamesboyle/golang-interview-prep/internal/middlewares"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/user"
)

func main() {
	dbURL := os.Getenv("db_uri")
	dbDriver := os.Getenv("db_driver")
	migrationPath := os.Getenv("migration_path")
	server_addr := os.Getenv("app_addr")

	err := run(server_addr, dbURL, dbDriver, migrationPath)

	if err != nil {
		log.Fatal(fmt.Errorf("main: run: %w", err))
	}
}

func run(serve_addr, dbURL, dbDriver, migrationPath string) error {
	db, err := dbInstance(dbURL, dbDriver)

	if err != nil {
		return fmt.Errorf("get_db_instance: %w", err)
	}

	defer db.Close()

	err = runMigrations(db, migrationPath)

	if err != nil {
		return fmt.Errorf("run_migration: %w", err)
	}

	svc := user.NewService(db)

	userHandler := user.Handler{Svc: *svc}

	http.HandleFunc("/user", middlewares.AuthRequired(userHandler.AddUser))

	log.Printf("starting http server on %s", serve_addr)
	err = http.ListenAndServe(serve_addr, nil)

	if err != nil {
		return fmt.Errorf("listen_on_addr: %w", err)
	}

	return nil
}

func dbInstance(dbURI, dbDriver string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbURI)
	if err != nil {
		return nil, fmt.Errorf("db_instance: %w", err)
	}
	return db, nil
}

func runMigrations(db *sql.DB, migrationPath string) error {
	// Create a new instance of the PostgreSQL driver for migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("get_instance: %w", err)
	}

	migrationFilePath := fmt.Sprintf("file://%s", migrationPath)
	m, err := migrate.NewWithDatabaseInstance(migrationFilePath, "", driver)
	if err != nil {
		return fmt.Errorf("create_migration_instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running_migration: %w", err)
	}

	fmt.Println("Database migration complete.")
	return nil
}
