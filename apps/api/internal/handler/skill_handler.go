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

type SkillHandler struct {
	svc service.SkillService
}

func NewSkillHandler(svc service.SkillService) *SkillHandler {
	return &SkillHandler{svc: svc}
}

func (h *SkillHandler) ListPublic(c *fiber.Ctx) error {
	skills, err := h.svc.ListSkills(c.Context(), true)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar skill",
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(skills, nil))
}

func (h *SkillHandler) ListAdmin(c *fiber.Ctx) error {
	skills, err := h.svc.ListSkills(c.Context(), false)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar skill",
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(skills, nil))
}

func (h *SkillHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	skill, err := h.svc.GetSkillByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrSkillNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Skill tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil skill",
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(skill, nil))
}

func (h *SkillHandler) Create(c *fiber.Ctx) error {
	var req domain.SkillCreateRequest
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
			"Data skill tidak valid",
			fields,
		))
	}

	skill, err := h.svc.CreateSkill(c.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrSkillConflict) {
			return c.Status(http.StatusConflict).JSON(domain.NewErrorResponse(
				"CONFLICT",
				err.Error(),
				nil,
			))
		}
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"CREATION_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusCreated).JSON(domain.NewSuccessResponse(skill, nil))
}

func (h *SkillHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.SkillUpdateRequest
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
			"Data skill tidak valid",
			fields,
		))
	}

	skill, err := h.svc.UpdateSkill(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrSkillNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Skill tidak ditemukan",
				nil,
			))
		}
		if errors.Is(err, repository.ErrSkillConflict) {
			return c.Status(http.StatusConflict).JSON(domain.NewErrorResponse(
				"CONFLICT",
				err.Error(),
				nil,
			))
		}
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"UPDATE_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(skill, nil))
}

func (h *SkillHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.svc.DeleteSkill(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrSkillNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Skill tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"DELETE_FAILED",
			"Gagal menghapus skill",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Skill berhasil dihapus",
	}, nil))
}
