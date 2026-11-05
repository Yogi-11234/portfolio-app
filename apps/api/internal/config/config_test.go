package config

import (
	"os"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	// Test short secret
	os.Setenv("JWT_SECRET", "too-short")
	os.Setenv("DATABASE_URL", "postgres://test")
	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for short JWT_SECRET, got nil")
	}

	// Test valid secret (>= 32 chars)
	os.Setenv("JWT_SECRET", "this-is-a-valid-secret-key-with-more-than-32-chars")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected valid config, got error: %v", err)
	}

	if cfg.AppPort != 8080 && cfg.AppPort != 0 {
		t.Logf("App port: %d", cfg.AppPort)
	}
}
