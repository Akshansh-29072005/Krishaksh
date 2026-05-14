package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aarcsx/krisho-backend/internal/config"
	"github.com/aarcsx/krisho-backend/internal/database"
)

func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer db.Close()

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "db/migrations"
	}

	fmt.Printf("Running migrations from %s\n", migrationsDir)
	if err := database.RunMigrations(context.Background(), db, migrationsDir); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Println("All migrations applied successfully.")
}
