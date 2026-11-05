package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"portfolio/api/internal/security"
)

func findEnvFile() string {
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

func main() {
	envPath := findEnvFile()
	_ = godotenv.Load(envPath)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	adminEmail := os.Getenv("ADMIN_SEED_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@local.test"
	}

	adminPassword := os.Getenv("ADMIN_SEED_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin123!Porto2026"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	var userID string
	err = pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", adminEmail).Scan(&userID)
	if err == nil {
		log.Printf("Admin user '%s' already exists (ID: %s).", adminEmail, userID)
		// Ensure profile exists
		ensureProfile(ctx, pool, userID, adminEmail)
		log.Println("Seeder completed successfully.")
		return
	}

	log.Printf("Seeding admin user '%s'...", adminEmail)
	passwordHash, err := security.HashPassword(adminPassword)
	if err != nil {
		log.Fatalf("Failed to hash password with argon2id: %v", err)
	}

	adminName := "Admin Portfolio"
	err = pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name, role)
		 VALUES ($1, $2, $3, 'admin')
		 RETURNING id`,
		adminEmail, passwordHash, adminName,
	).Scan(&userID)
	if err != nil {
		log.Fatalf("Failed to insert admin user: %v", err)
	}

	log.Printf("Admin user created with ID: %s", userID)
	ensureProfile(ctx, pool, userID, adminEmail)
	log.Println("Seeder completed successfully.")
}

func ensureProfile(ctx context.Context, pool *pgxpool.Pool, userID, email string) {
	_, err := pool.Exec(ctx,
		`INSERT INTO profile (user_id, full_name, headline, bio_md, bio_html, location, public_email, open_to_work)
		 VALUES ($1, 'Admin Portfolio', 'Senior Fullstack Engineer', 'Senior fullstack engineer building modern web applications.', '<p>Senior fullstack engineer building modern web applications.</p>', 'Jakarta, Indonesia', $2, true)
		 ON CONFLICT (user_id) DO NOTHING`,
		userID, email,
	)
	if err != nil {
		log.Fatalf("Failed to ensure profile: %v", err)
	}
	log.Println("Profile record ensured.")
}
