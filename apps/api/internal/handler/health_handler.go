package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

type HealthHandler struct {
	pool *pgxpool.Pool
}

func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "connected"
	if err := h.pool.Ping(ctx); err != nil {
		dbStatus = "disconnected"
	}

	data := fiber.Map{
		"status":    "healthy",
		"database":  dbStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(data, nil))
}
