package middleware

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/ratelimit"
)

// RateLimiter is a re-export alias so existing callers (handler, server) don't
// need to change their import paths. All types satisfy ratelimit.Limiter.
type RateLimiter = ratelimit.Limiter

// NewMemoryRateLimiter creates an in-memory sliding-window RateLimiter.
// It is the primary constructor used in main.go and tests.
func NewMemoryRateLimiter() RateLimiter {
	return ratelimit.NewMemoryLimiter()
}

// GlobalRateLimit enforces a per-IP request rate limit (e.g. 100 req/minute).
func GlobalRateLimit(limiter RateLimiter, limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		key := "global:" + ip
		allowed, retryAfter := limiter.Allow(key, limit, window)
		if !allowed {
			c.Set("Retry-After", ratelimit.FormatRetryAfter(retryAfter))
			return c.Status(http.StatusTooManyRequests).JSON(domain.NewErrorResponse(
				"RATE_LIMIT_EXCEEDED",
				"Terlalu banyak permintaan. Silakan tunggu beberapa saat lagi.",
				nil,
			))
		}
		return c.Next()
	}
}
