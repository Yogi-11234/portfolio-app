package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/service"
	"portfolio/api/internal/validator"
)

type MessageHandler struct {
	svc service.MessageService
}

func NewMessageHandler(svc service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

// CreateMessage handles POST /api/v1/contact — public, no auth required.
func (h *MessageHandler) CreateMessage(c *fiber.Ctx) error {
	var req domain.MessageCreateRequest
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
			"Data pesan tidak valid",
			fields,
		))
	}

	ip := c.IP()
	userAgent := c.Get("User-Agent")

	_, err := h.svc.CreateMessage(c.Context(), req, ip, userAgent)
	if err != nil {
		// Honeypot: silently respond with 200 to avoid leaking detection logic.
		if errors.Is(err, service.ErrHoneypotTriggered) {
			return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
				"message": "Pesan berhasil terkirim. Terima kasih!",
			}, nil))
		}
		if errors.Is(err, service.ErrContactRateLimited) {
			return c.Status(http.StatusTooManyRequests).JSON(domain.NewErrorResponse(
				"RATE_LIMIT_EXCEEDED",
				"Terlalu banyak pesan dikirim. Silakan coba lagi nanti.",
				nil,
			))
		}
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"SEND_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusCreated).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Pesan berhasil terkirim. Terima kasih!",
	}, nil))
}

// ListMessages handles GET /api/v1/admin/messages — admin only.
func (h *MessageHandler) ListMessages(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	status := c.Query("status")

	messages, total, err := h.svc.ListMessages(c.Context(), status, page, perPage)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar pesan",
			nil,
		))
	}

	meta := &domain.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(messages, meta))
}

// GetMessageByID handles GET /api/v1/admin/messages/:id — admin only.
func (h *MessageHandler) GetMessageByID(c *fiber.Ctx) error {
	id := c.Params("id")
	msg, err := h.svc.GetMessageByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrMessageNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Pesan tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil pesan",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(msg, nil))
}

// UpdateStatus handles PATCH /api/v1/admin/messages/:id/status — admin + CSRF.
func (h *MessageHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.MessageStatusUpdateRequest
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
			"Status tidak valid",
			fields,
		))
	}

	if err := h.svc.UpdateMessageStatus(c.Context(), id, req.Status); err != nil {
		if errors.Is(err, repository.ErrMessageNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Pesan tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"UPDATE_FAILED",
			"Gagal memperbarui status pesan",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Status pesan berhasil diperbarui",
	}, nil))
}

// DeleteMessage handles DELETE /api/v1/admin/messages/:id — admin + CSRF.
func (h *MessageHandler) DeleteMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteMessage(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrMessageNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Pesan tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"DELETE_FAILED",
			"Gagal menghapus pesan",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Pesan berhasil dihapus",
	}, nil))
}
