package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func findEnvFile() string {
	// Check current dir, then parent dirs
	dir, err := os.Getwd()
	if err != nil {
		return ".env"
	}
	for i := 0; i < 4; i++ {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ".env"
}

func findMigrationsDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "migrations"
	}
	for i := 0; i < 4; i++ {
		candidate := filepath.Join(dir, "migrations")
		if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
			return candidate
		}
		candidate2 := filepath.Join(dir, "apps", "api", "migrations")
		if stat, err := os.Stat(candidate2); err == nil && stat.IsDir() {
			return candidate2
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "migrations"
}

func main() {
	envPath := findEnvFile()
	_ = godotenv.Load(envPath)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	migrationsDir := findMigrationsDir()
	sourceURL := fmt.Sprintf("file://%s", filepath.ToSlash(migrationsDir))

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}
	defer m.Close()

	flag.Parse()
	action := "up"
	if len(flag.Args()) > 0 {
		action = strings.ToLower(flag.Args()[0])
	}

	switch action {
	case "up":
		log.Printf("Applying migrations from %s...", sourceURL)
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No migration changes to apply (database is already up to date).")
			} else {
				log.Fatalf("Migration Up failed: %v", err)
			}
		} else {
			log.Println("Migration Up applied successfully.")
		}
	case "down":
		log.Printf("Rolling back 1 migration step...")
		if err := m.Steps(-1); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No migrations to rollback.")
			} else {
				log.Fatalf("Migration Down failed: %v", err)
			}
		} else {
			log.Println("Migration Down rolled back successfully.")
		}
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		log.Printf("Current migration version: %d (dirty: %v)", v, dirty)
	default:
		log.Fatalf("Unknown migration action: %s. Use 'up', 'down', or 'version'.", action)
	}
}
