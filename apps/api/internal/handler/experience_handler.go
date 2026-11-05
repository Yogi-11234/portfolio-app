package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/service"
	"portfolio/api/internal/validator"
)

type ExperienceHandler struct {
	svc service.ExperienceService
}

func NewExperienceHandler(svc service.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{svc: svc}
}

func (h *ExperienceHandler) ListPublic(c *fiber.Ctx) error {
	expType := c.Query("type")
	experiences, err := h.svc.ListExperiences(c.Context(), true, expType)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"INVALID_QUERY",
			err.Error(),
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(experiences, nil))
}

func (h *ExperienceHandler) ListAdmin(c *fiber.Ctx) error {
	expType := c.Query("type")
	experiences, err := h.svc.ListExperiences(c.Context(), false, expType)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"INVALID_QUERY",
			err.Error(),
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(experiences, nil))
}

func (h *ExperienceHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	experience, err := h.svc.GetExperienceByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrExperienceNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Pengalaman tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil pengalaman",
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(experience, nil))
}

func (h *ExperienceHandler) Create(c *fiber.Ctx) error {
	var req domain.ExperienceCreateRequest
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
			"Data pengalaman tidak valid",
			fields,
		))
	}

	exp, err := h.svc.CreateExperience(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"CREATION_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusCreated).JSON(domain.NewSuccessResponse(exp, nil))
}

func (h *ExperienceHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.ExperienceUpdateRequest
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
			"Data pengalaman tidak valid",
			fields,
		))
	}

	exp, err := h.svc.UpdateExperience(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrExperienceNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Pengalaman tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"UPDATE_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(exp, nil))
}

func (h *ExperienceHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteExperience(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrExperienceNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Pengalaman tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"DELETE_FAILED",
			"Gagal menghapus pengalaman",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Pengalaman berhasil dihapus",
	}, nil))
}
