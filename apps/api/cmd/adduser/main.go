package main

import (
	"context"
	"fmt"
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
	_ = godotenv.Load(findEnvFile())

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	targetEmail := "yogi.aditia2004@gmail.com"
	targetPassword := "lupa_lagi*1125"
	targetName := "Yogi Aditia"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// List current users
	rows, err := pool.Query(ctx, "SELECT id, email, name, role FROM users")
	if err == nil {
		fmt.Println("--- Current Users in Database ---")
		for rows.Next() {
			var id, email, name, role string
			if err := rows.Scan(&id, &email, &name, &role); err == nil {
				fmt.Printf("User: ID=%s, Email=%s, Name=%s, Role=%s\n", id, email, name, role)
			}
		}
		rows.Close()
	}

	// Hash the password with argon2id
	hash, err := security.HashPassword(targetPassword)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	var userID string
	err = pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", targetEmail).Scan(&userID)
	if err == nil {
		// User already exists, update password and reset failed logins / locks
		_, err = pool.Exec(ctx, `
			UPDATE users 
			SET password_hash = $1, name = $2, role = 'admin', failed_logins = 0, locked_until = NULL, updated_at = now()
			WHERE id = $3
		`, hash, targetName, userID)
		if err != nil {
			log.Fatalf("Failed to update user: %v", err)
		}
		fmt.Printf("Successfully updated password and details for user: %s (ID: %s)\n", targetEmail, userID)
	} else {
		// Insert new user
		err = pool.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, name, role)
			VALUES ($1, $2, $3, 'admin')
			RETURNING id
		`, targetEmail, hash, targetName).Scan(&userID)
		if err != nil {
			log.Fatalf("Failed to insert user: %v", err)
		}
		fmt.Printf("Successfully created new admin user: %s (ID: %s)\n", targetEmail, userID)
	}

	// Check if profile exists for this user
	var profileID string
	err = pool.QueryRow(ctx, "SELECT id FROM profile WHERE user_id = $1", userID).Scan(&profileID)
	if err != nil {
		// Check if there is an existing profile from another user that we should copy or update
		var existingProfID string
		err = pool.QueryRow(ctx, "SELECT id FROM profile LIMIT 1").Scan(&existingProfID)
		if err == nil {
			// Update the existing profile to belong to this new user
			_, err = pool.Exec(ctx, "UPDATE profile SET user_id = $1, full_name = $2, public_email = $3 WHERE id = $4",
				userID, targetName, targetEmail, existingProfID)
			if err != nil {
				// If error (e.g. unique constraint), insert new
				_, _ = pool.Exec(ctx, `
					INSERT INTO profile (user_id, full_name, headline, bio_md, bio_html, location, public_email, open_to_work)
					VALUES ($1, $2, 'Fullstack Software Engineer', '', '', 'Jakarta, Indonesia', $3, true)
					ON CONFLICT (user_id) DO NOTHING
				`, userID, targetName, targetEmail)
			} else {
				fmt.Printf("Linked existing profile (%s) to user %s\n", existingProfID, targetEmail)
			}
		} else {
			_, _ = pool.Exec(ctx, `
				INSERT INTO profile (user_id, full_name, headline, bio_md, bio_html, location, public_email, open_to_work)
				VALUES ($1, $2, 'Fullstack Software Engineer', '', '', 'Jakarta, Indonesia', $3, true)
				ON CONFLICT (user_id) DO NOTHING
			`, userID, targetName, targetEmail)
			fmt.Printf("Created new profile for user %s\n", targetEmail)
		}
	} else {
		fmt.Printf("Profile already exists for user %s (ID: %s)\n", targetEmail, profileID)
	}

	fmt.Println("All done!")
}
