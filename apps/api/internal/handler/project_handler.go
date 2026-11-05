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

type ProjectHandler struct {
	svc service.ProjectService
}

func NewProjectHandler(svc service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) ListPublic(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	featuredStr := c.Query("featured")
	var featuredOnly *bool
	if featuredStr == "true" {
		t := true
		featuredOnly = &t
	}

	filter := domain.ProjectFilter{
		CategorySlug: c.Query("category"),
		SkillSlug:    c.Query("skill"),
		FeaturedOnly: featuredOnly,
		Status:       "published", // Always published for public
		Page:         page,
		PerPage:      12,
	}

	projects, total, err := h.svc.ListProjects(c.Context(), filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar proyek",
			nil,
		))
	}

	meta := &domain.Meta{
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Total:   total,
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(projects, meta))
}

func (h *ProjectHandler) GetBySlugPublic(c *fiber.Ctx) error {
	slug := c.Params("slug")
	project, err := h.svc.GetProjectBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Proyek tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil detail proyek",
			nil,
		))
	}

	if project.Status != "published" {
		return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
			"NOT_FOUND",
			"Proyek tidak ditemukan",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(project, nil))
}

func (h *ProjectHandler) ListAdmin(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	status := c.Query("status")

	filter := domain.ProjectFilter{
		CategorySlug: c.Query("category"),
		Status:       status, // draft, published, archived or all
		Page:         page,
		PerPage:      perPage,
	}

	projects, total, err := h.svc.ListProjects(c.Context(), filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar proyek",
			nil,
		))
	}

	meta := &domain.Meta{
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Total:   total,
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(projects, meta))
}

func (h *ProjectHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	project, err := h.svc.GetProjectByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Proyek tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil proyek",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(project, nil))
}

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	var req domain.ProjectCreateRequest
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
			"Data proyek tidak valid",
			fields,
		))
	}

	project, err := h.svc.CreateProject(c.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrProjectConflict) {
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

	return c.Status(http.StatusCreated).JSON(domain.NewSuccessResponse(project, nil))
}

func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.ProjectUpdateRequest
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
			"Data proyek tidak valid",
			fields,
		))
	}

	project, err := h.svc.UpdateProject(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Proyek tidak ditemukan",
				nil,
			))
		}
		if errors.Is(err, repository.ErrProjectConflict) {
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

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(project, nil))
}

func (h *ProjectHandler) UpdateSkills(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.ProjectSkillsUpdateRequest
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
			"Data skill_ids tidak valid",
			fields,
		))
	}

	if err := h.svc.UpdateProjectSkills(c.Context(), id, req.SkillIDs); err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Proyek tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"UPDATE_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Skills proyek berhasil diperbarui",
	}, nil))
}

func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteProject(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Proyek tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"DELETE_FAILED",
			"Gagal menghapus proyek",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Proyek berhasil dihapus",
	}, nil))
}
