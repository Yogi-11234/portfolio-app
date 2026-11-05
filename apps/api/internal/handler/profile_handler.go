package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/service"
	"portfolio/api/internal/validator"
)

type ProfileHandler struct {
	svc service.ProfileService
}

func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) Get(c *fiber.Ctx) error {
	profile, err := h.svc.GetProfile(c.Context())
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
			"NOT_FOUND",
			"Profil belum dibuat",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(profile, nil))
}

func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	userVal := c.Locals("user")
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.NewErrorResponse(
			"UNAUTHORIZED",
			"Otorisasi diperlukan",
			nil,
		))
	}

	var req domain.ProfileUpdateRequest
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
			"Validasi formulir profil gagal",
			fields,
		))
	}

	updatedProfile, err := h.svc.UpdateProfile(c.Context(), user.ID, req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"OPERATION_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(updatedProfile, nil))
}
