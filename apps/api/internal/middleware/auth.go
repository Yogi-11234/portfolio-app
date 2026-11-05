package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/service"
)

// RequireAuth validates the JWT in the httpOnly auth_token cookie.
func RequireAuth(authSvc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("auth_token")
		if cookie == "" {
			return c.Status(http.StatusUnauthorized).JSON(domain.NewErrorResponse(
				"UNAUTHORIZED",
				"Sesi tidak ditemukan atau telah berakhir. Silakan login kembali.",
				nil,
			))
		}

		claims, err := authSvc.ValidateToken(cookie)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(domain.NewErrorResponse(
				"UNAUTHORIZED",
				"Token sesi tidak valid atau telah kedaluwarsa.",
				nil,
			))
		}

		user, err := authSvc.GetUserByID(c.Context(), claims.UserID)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(domain.NewErrorResponse(
				"UNAUTHORIZED",
				"Pengguna tidak ditemukan.",
				nil,
			))
		}

		if user.Role != "admin" {
			return c.Status(http.StatusForbidden).JSON(domain.NewErrorResponse(
				"FORBIDDEN",
				"Akses ditolak: role tidak memiliki izin.",
				nil,
			))
		}

		c.Locals("user", user)
		return c.Next()
	}
}
