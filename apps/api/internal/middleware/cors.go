package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS configures CORS policies based on the allowed origins whitelist.
func CORS(allowedOrigins []string) fiber.Handler {
	originsStr := strings.Join(allowedOrigins, ", ")
	return cors.New(cors.Config{
		AllowOrigins:     originsStr,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID,X-CSRF-Token",
		ExposeHeaders:    "X-Request-ID,X-CSRF-Token",
	})
}
