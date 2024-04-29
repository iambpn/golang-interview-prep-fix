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

	"github.com/matthewjamesboyle/golang-interview-prep/internal/user"
)

func main() {
	dbURL := os.Getenv("db_uri")
	dbType := os.Getenv("db_type")
	migrationPath := os.Getenv("migration_path")
	server_addr := os.Getenv("app_addr")

	err := run(server_addr, dbURL, dbType, migrationPath)

	if err != nil {
		log.Fatal(fmt.Errorf("main: run: %w", err))
	}
}

func run(serve_addr, dbURL, dbType, migrationPath string) error {
	err := runMigrations(dbURL, dbType, migrationPath)

	if err != nil {
		return fmt.Errorf("run_migration: %w", err)
	}

	svc, err := user.NewService("admin", "admin")
	if err != nil {
		return fmt.Errorf("initialize_user_service: %w", err)
	}

	h := user.Handler{Svc: *svc}

	http.HandleFunc("/user", h.AddUser)

	log.Printf("starting http server on %s", serve_addr)
	err = http.ListenAndServe(serve_addr, nil)

	if err != nil {
		return fmt.Errorf("listen_on_addr: %w", err)
	}

	return nil
}

func runMigrations(dbURL, dbType, migrationPath string) error {
	db, err := sql.Open(dbType, dbURL)
	if err != nil {
		return fmt.Errorf("connect_to_db: %w", err)
	}
	defer db.Close()

	// Create a new instance of the PostgreSQL driver for migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("get_instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationPath), dbType, driver)
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
