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

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) List(c *fiber.Ctx) error {
	catType := c.Query("type")
	categories, err := h.svc.ListCategories(c.Context(), catType)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"INVALID_QUERY",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(categories, nil))
}

func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	category, err := h.svc.GetCategoryByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Kategori tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil kategori",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(category, nil))
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var req domain.CategoryCreateRequest
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
			"Data kategori tidak valid",
			fields,
		))
	}

	category, err := h.svc.CreateCategory(c.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryConflict) {
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

	return c.Status(http.StatusCreated).JSON(domain.NewSuccessResponse(category, nil))
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.CategoryUpdateRequest
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
			"Data kategori tidak valid",
			fields,
		))
	}

	category, err := h.svc.UpdateCategory(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Kategori tidak ditemukan",
				nil,
			))
		}
		if errors.Is(err, repository.ErrCategoryConflict) {
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

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(category, nil))
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.svc.DeleteCategory(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Kategori tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"DELETE_FAILED",
			"Gagal menghapus kategori",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Kategori berhasil dihapus",
	}, nil))
}
