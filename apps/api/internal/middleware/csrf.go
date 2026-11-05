package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/domain"
)

// ValidateCSRF enforces double-submit CSRF token on mutating requests (POST, PUT, PATCH, DELETE).
func ValidateCSRF() fiber.Handler {
	return func(c *fiber.Ctx) error {
		method := c.Method()
		if method == fiber.MethodGet || method == fiber.MethodHead || method == fiber.MethodOptions {
			return c.Next()
		}

		cookieToken := c.Cookies("csrf_token")
		headerToken := c.Get("X-CSRF-Token")

		if cookieToken == "" || headerToken == "" {
			return c.Status(http.StatusForbidden).JSON(domain.NewErrorResponse(
				"CSRF_TOKEN_MISSING",
				"Token CSRF tidak ditemukan. Silakan refresh halaman dan coba lagi.",
				nil,
			))
		}

		if subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
			return c.Status(http.StatusForbidden).JSON(domain.NewErrorResponse(
				"CSRF_TOKEN_INVALID",
				"Token CSRF tidak cocok. Akses ditolak demi keamanan.",
				nil,
			))
		}

		return c.Next()
	}
}
