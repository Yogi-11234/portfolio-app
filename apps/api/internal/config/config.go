package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                   string
	AppPort                  int
	AppBaseURL               string
	DatabaseURL              string
	JWTSecret                string
	JWTTLHours               int
	CookieSecure             bool
	CORSAllowedOrigins       []string
	StorageDriver            string
	StorageLocalDir          string
	UploadMaxBytes           int64
	RateLimitContactPerHour  int
	AdminSeedEmail           string
	AdminSeedPassword        string
	SupabaseURL              string
	SupabaseBucket           string
	SupabaseKey              string
}

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

// Load loads the configuration from environment variables or .env file.
func Load() (*Config, error) {
	envPath := findEnvFile()
	_ = godotenv.Load(envPath)

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(strings.TrimSpace(jwtSecret)) < 32 {
		return nil, errors.New("FATAL: JWT_SECRET must be at least 32 characters long and not empty")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("FATAL: DATABASE_URL is required")
	}

	appPort := 8080
	if portStr := os.Getenv("APP_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil && p > 0 {
			appPort = p
		}
	}

	jwtTTL := 12
	if ttlStr := os.Getenv("JWT_TTL_HOURS"); ttlStr != "" {
		if t, err := strconv.Atoi(ttlStr); err == nil && t > 0 {
			jwtTTL = t
		}
	}

	cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
	if os.Getenv("APP_ENV") == "production" {
		cookieSecure = true
	}

	corsOriginsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	var corsOrigins []string
	if corsOriginsStr != "" {
		for _, o := range strings.Split(corsOriginsStr, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				corsOrigins = append(corsOrigins, trimmed)
			}
		}
	}
	if len(corsOrigins) == 0 {
		corsOrigins = []string{"http://localhost:5173"}
	}

	uploadMax := int64(2097152) // 2 MB default
	if maxStr := os.Getenv("UPLOAD_MAX_BYTES"); maxStr != "" {
		if m, err := strconv.ParseInt(maxStr, 10, 64); err == nil && m > 0 {
			uploadMax = m
		}
	}

	contactRateLimit := 5
	if rlStr := os.Getenv("RATE_LIMIT_CONTACT_PER_HOUR"); rlStr != "" {
		if r, err := strconv.Atoi(rlStr); err == nil && r > 0 {
			contactRateLimit = r
		}
	}

	storageLocalDir := os.Getenv("STORAGE_LOCAL_DIR")
	if storageLocalDir == "" {
		storageLocalDir = "./storage/uploads"
	}

	storageDriver := strings.ToLower(strings.TrimSpace(os.Getenv("STORAGE_DRIVER")))
	if storageDriver == "" {
		storageDriver = "supabase"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	supabaseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("SUPABASE_URL")), "/")
	if supabaseURL == "" {
		supabaseURL = "https://kpbuhrpgnrqmwrzzxols.supabase.co"
	}

	supabaseBucket := strings.TrimSpace(os.Getenv("SUPABASE_STORAGE_BUCKET"))
	if supabaseBucket == "" {
		supabaseBucket = strings.TrimSpace(os.Getenv("SUPABASE_BUCKET"))
	}
	if supabaseBucket == "" {
		supabaseBucket = "portfolio_assets"
	}

	supabaseKey := strings.TrimSpace(os.Getenv("SUPABASE_KEY"))
	if supabaseKey == "" {
		supabaseKey = strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	}
	if supabaseKey == "" {
		supabaseKey = strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY"))
	}

	appBaseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("APP_BASE_URL")), "/")
	if appBaseURL == "" {
		if storageDriver == "supabase" {
			appBaseURL = fmt.Sprintf("%s/storage/v1/object/public/%s", supabaseURL, supabaseBucket)
		} else {
			appBaseURL = fmt.Sprintf("http://localhost:%d", appPort)
		}
	}

	return &Config{
		AppEnv:                  appEnv,
		AppPort:                 appPort,
		AppBaseURL:              appBaseURL,
		DatabaseURL:             databaseURL,
		JWTSecret:               jwtSecret,
		JWTTLHours:              jwtTTL,
		CookieSecure:            cookieSecure,
		CORSAllowedOrigins:      corsOrigins,
		StorageDriver:           storageDriver,
		StorageLocalDir:         storageLocalDir,
		UploadMaxBytes:          uploadMax,
		RateLimitContactPerHour: contactRateLimit,
		AdminSeedEmail:          os.Getenv("ADMIN_SEED_EMAIL"),
		AdminSeedPassword:       os.Getenv("ADMIN_SEED_PASSWORD"),
		SupabaseURL:             supabaseURL,
		SupabaseBucket:          supabaseBucket,
		SupabaseKey:             supabaseKey,
	}, nil
}
