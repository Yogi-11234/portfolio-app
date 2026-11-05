package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/domain"
)

// Recover catches any panic, logs the stack trace internally, and returns a generic 500 error envelope.
func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				reqID := c.Locals("request_id")
				slog.Error("Unhandled panic recovered",
					"error", r,
					"request_id", reqID,
					"path", c.Path(),
					"method", c.Method(),
					"stack", string(debug.Stack()),
				)

				_ = c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
					"INTERNAL_SERVER_ERROR",
					"Terjadi kesalahan internal pada server",
					nil,
				))
			}
		}()
		return c.Next()
	}
}
