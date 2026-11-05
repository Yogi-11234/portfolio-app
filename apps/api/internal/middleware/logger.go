package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDAndLogger generates a UUID request_id and outputs structured logs without sensitive data.
func RequestIDAndLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Locals("request_id", reqID)
		c.Set("X-Request-ID", reqID)

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		status := c.Response().StatusCode()

		// Structured log: method, path, status, duration, request_id, IP
		// NEVER log body, password, token, or cookie!
		slog.Info("HTTP Request",
			"request_id", reqID,
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"ip", c.IP(),
		)

		return err
	}
}
