package service

import (
	"context"
	"testing"
	"time"

	"portfolio/api/internal/config"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/security"
)

type mockUserRepo struct {
	user        *domain.User
	failedCalls int
	successCall bool
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.user != nil && m.user.Email == email {
		return m.user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.user != nil && m.user.ID == id {
		return m.user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepo) RecordFailedLogin(ctx context.Context, id string, lockDuration time.Duration) error {
	m.failedCalls++
	return nil
}

func (m *mockUserRepo) RecordSuccessfulLogin(ctx context.Context, id string) error {
	m.successCall = true
	return nil
}

func TestAuthService_LoginAndTokenValidation(t *testing.T) {
	pwdHash, _ := security.HashPassword("CorrectPassword123!")
	mockRepo := &mockUserRepo{
		user: &domain.User{
			ID:           "test-uuid-123",
			Email:        "admin@test.local",
			PasswordHash: pwdHash,
			Name:         "Admin Test",
			Role:         "admin",
		},
	}

	cfg := &config.Config{
		JWTSecret:  "super-secret-key-that-is-at-least-32-chars-long!",
		JWTTLHours: 12,
	}

	authSvc := NewAuthService(mockRepo, cfg)
	ctx := context.Background()

	// 1. Test successful login
	user, token, csrf, err := authSvc.Login(ctx, domain.LoginRequest{
		Email:    "admin@test.local",
		Password: "CorrectPassword123!",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if user.ID != "test-uuid-123" {
		t.Errorf("Expected user ID test-uuid-123, got %s", user.ID)
	}
	if token == "" || csrf == "" {
		t.Error("Expected non-empty token and csrf")
	}

	// 2. Validate token
	claims, err := authSvc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.UserID != "test-uuid-123" || claims.Role != "admin" {
		t.Errorf("Unexpected claims: %+v", claims)
	}

	// 3. Test wrong password
	_, _, _, err = authSvc.Login(ctx, domain.LoginRequest{
		Email:    "admin@test.local",
		Password: "WrongPassword123!",
	})
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
	if mockRepo.failedCalls != 1 {
		t.Errorf("Expected 1 failed call, got %d", mockRepo.failedCalls)
	}

	// 4. Test locked account
	lockTime := time.Now().Add(10 * time.Minute)
	mockRepo.user.LockedUntil = &lockTime
	_, _, _, err = authSvc.Login(ctx, domain.LoginRequest{
		Email:    "admin@test.local",
		Password: "CorrectPassword123!",
	})
	if err != ErrAccountLocked {
		t.Errorf("Expected ErrAccountLocked, got %v", err)
	}
}
