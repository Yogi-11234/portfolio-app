package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/config"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/middleware"
	"portfolio/api/internal/service"
	"portfolio/api/internal/validator"
)

type AuthHandler struct {
	authSvc     service.AuthService
	cfg         *config.Config
	rateLimiter middleware.RateLimiter
}

func NewAuthHandler(authSvc service.AuthService, cfg *config.Config, rateLimiter middleware.RateLimiter) *AuthHandler {
	return &AuthHandler{
		authSvc:     authSvc,
		cfg:         cfg,
		rateLimiter: rateLimiter,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"INVALID_REQUEST",
			"Format payload JSON tidak valid",
			nil,
		))
	}

	if fields := validator.ValidateStruct(req); len(fields) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(domain.NewErrorResponse(
			"VALIDATION_ERROR",
			"Data login tidak memenuhi validasi",
			fields,
		))
	}

	// Rate limit: 5 requests per 15 minutes per IP + Email
	cleanEmail := strings.ToLower(strings.TrimSpace(req.Email))
	rateLimitKey := fmt.Sprintf("login:%s:%s", c.IP(), cleanEmail)
	allowed, retryAfter := h.rateLimiter.Allow(rateLimitKey, 5, 15*time.Minute)
	if !allowed {
		c.Set("Retry-After", fmt.Sprintf("%d", int(retryAfter.Seconds())+1))
		return c.Status(http.StatusTooManyRequests).JSON(domain.NewErrorResponse(
			"TOO_MANY_ATTEMPTS",
			"Terlalu banyak percobaan login gagal. Silakan coba kembali nanti.",
			nil,
		))
	}

	user, token, csrfToken, err := h.authSvc.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrAccountLocked) {
			return c.Status(http.StatusTooManyRequests).JSON(domain.NewErrorResponse(
				"ACCOUNT_LOCKED",
				service.ErrAccountLocked.Error(),
				nil,
			))
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(http.StatusUnauthorized).JSON(domain.NewErrorResponse(
				"INVALID_CREDENTIALS",
				"Email atau password yang Anda masukkan salah.",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal memproses login",
			nil,
		))
	}

	maxAgeSeconds := h.cfg.JWTTLHours * 3600

	// Set auth_token (httpOnly, Strict, Secure in prod)
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAgeSeconds,
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: "Strict",
	})

	// Set csrf_token (non-httpOnly for JavaScript reading)
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		MaxAge:   maxAgeSeconds,
		HTTPOnly: false,
		Secure:   h.cfg.CookieSecure,
		SameSite: "Strict",
	})

	responseData := fiber.Map{
		"user":       user.ToResponse(),
		"csrf_token": csrfToken,
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(responseData, nil))
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Expire auth_token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: "Strict",
	})

	// Expire csrf_token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: false,
		Secure:   h.cfg.CookieSecure,
		SameSite: "Strict",
	})

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Berhasil logout",
	}, nil))
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userVal := c.Locals("user")
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.NewErrorResponse(
			"UNAUTHORIZED",
			"Pengguna tidak terautentikasi",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(user.ToResponse(), nil))
}
