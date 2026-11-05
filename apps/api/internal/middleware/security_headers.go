package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders applies production security headers according to PRD Bagian 4.5
func SecurityHeaders(isSecure bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Base CSP
		csp := "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self'; frame-ancestors 'none'; object-src 'none'"

		// For PDF viewer on /uploads/*, allow frame-src 'self'
		if strings.HasPrefix(path, "/uploads") {
			csp = "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self'; frame-ancestors 'none'; frame-src 'self'; object-src 'self'"
		}

		c.Set("Content-Security-Policy", csp)
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		if isSecure {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		return c.Next()
	}
}
