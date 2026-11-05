package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/config"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/middleware"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/security"
	"portfolio/api/internal/service"
)

type mockUserRepoForHandler struct {
	user        *domain.User
	failedCalls int
	successCall bool
}

func (m *mockUserRepoForHandler) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.user != nil && m.user.Email == email {
		return m.user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepoForHandler) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.user != nil && m.user.ID == id {
		return m.user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepoForHandler) RecordFailedLogin(ctx context.Context, id string, lockDuration time.Duration) error {
	m.failedCalls++
	return nil
}

func (m *mockUserRepoForHandler) RecordSuccessfulLogin(ctx context.Context, id string) error {
	m.successCall = true
	return nil
}

func setupTestApp() (*fiber.App, *mockUserRepoForHandler, service.AuthService) {
	pwdHash, _ := security.HashPassword("CorrectPassword123!")
	mockRepo := &mockUserRepoForHandler{
		user: &domain.User{
			ID:           "test-admin-uuid",
			Email:        "admin@test.local",
			PasswordHash: pwdHash,
			Name:         "Admin Test",
			Role:         "admin",
		},
	}

	cfg := &config.Config{
		JWTSecret:          "secret-key-that-is-at-least-32-chars-long!",
		JWTTLHours:         12,
		CookieSecure:       false,
		CORSAllowedOrigins: []string{"http://localhost:5173"},
	}

	rateLimiter := middleware.NewMemoryRateLimiter()
	authSvc := service.NewAuthService(mockRepo, cfg)
	authHandler := NewAuthHandler(authSvc, cfg, rateLimiter)

	app := fiber.New()
	app.Use(middleware.Recover())
	app.Use(middleware.RequestIDAndLogger())
	app.Use(middleware.SecurityHeaders(false))

	api := app.Group("/api/v1")
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/logout", authHandler.Logout)
	api.Get("/auth/me", middleware.RequireAuth(authSvc), authHandler.Me)

	admin := api.Group("/admin", middleware.RequireAuth(authSvc), middleware.ValidateCSRF())
	admin.Post("/test-mutation", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	return app, mockRepo, authSvc
}

func TestAuthHandler_LoginFlow(t *testing.T) {
	app, _, _ := setupTestApp()

	// 1. Invalid payload format
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 for bad json, got %d", resp.StatusCode)
	}

	// 2. Validation error (short password)
	body, _ := json.Marshal(domain.LoginRequest{
		Email:    "admin@test.local",
		Password: "short",
	})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422 for validation error, got %d", resp.StatusCode)
	}

	// 3. Wrong password
	body, _ = json.Marshal(domain.LoginRequest{
		Email:    "admin@test.local",
		Password: "WrongPassword123!",
	})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong credentials, got %d", resp.StatusCode)
	}

	// 4. Successful login
	body, _ = json.Marshal(domain.LoginRequest{
		Email:    "admin@test.local",
		Password: "CorrectPassword123!",
	})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 for correct login, got %d", resp.StatusCode)
	}

	var authCookie *http.Cookie
	var csrfCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth_token" {
			authCookie = cookie
		}
		if cookie.Name == "csrf_token" {
			csrfCookie = cookie
		}
	}

	if authCookie == nil || authCookie.Value == "" || !authCookie.HttpOnly {
		t.Fatal("Expected httpOnly auth_token cookie")
	}
	if csrfCookie == nil || csrfCookie.Value == "" || csrfCookie.HttpOnly {
		t.Fatal("Expected non-httpOnly csrf_token cookie")
	}

	// 5. Test /auth/me with auth_token cookie
	req = httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.AddCookie(authCookie)
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for /auth/me with cookie, got %d", resp.StatusCode)
	}

	// 6. Test /auth/me without cookie (must be 401)
	req = httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 for /auth/me without cookie, got %d", resp.StatusCode)
	}

	// 7. Test Admin Mutation without CSRF token (must be 403)
	req = httptest.NewRequest("POST", "/api/v1/admin/test-mutation", nil)
	req.AddCookie(authCookie)
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 for admin mutation without CSRF, got %d", resp.StatusCode)
	}

	// 8. Test Admin Mutation with valid CSRF double-submit (must be 200)
	req = httptest.NewRequest("POST", "/api/v1/admin/test-mutation", nil)
	req.AddCookie(authCookie)
	req.AddCookie(csrfCookie)
	req.Header.Set("X-CSRF-Token", csrfCookie.Value)
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for admin mutation with valid CSRF, got %d", resp.StatusCode)
	}

	// 9. Test Logout (expires cookies)
	req = httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for logout, got %d", resp.StatusCode)
	}
}
