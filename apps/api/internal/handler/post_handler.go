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

type PostHandler struct {
	svc service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) ListPublic(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	filter := domain.PostFilter{
		CategorySlug: c.Query("category"),
		Status:       "published",
		Page:         page,
		PerPage:      10,
	}

	posts, total, err := h.svc.ListPosts(c.Context(), filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar postingan",
			nil,
		))
	}

	meta := &domain.Meta{
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Total:   total,
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(posts, meta))
}

func (h *PostHandler) GetBySlugPublic(c *fiber.Ctx) error {
	slug := c.Params("slug")
	post, err := h.svc.GetPostBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Postingan tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil detail postingan",
			nil,
		))
	}

	if post.Status != "published" {
		return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
			"NOT_FOUND",
			"Postingan tidak ditemukan",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(post, nil))
}

func (h *PostHandler) ListAdmin(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	filter := domain.PostFilter{
		CategorySlug: c.Query("category"),
		Status:       c.Query("status"),
		Page:         page,
		PerPage:      perPage,
	}

	posts, total, err := h.svc.ListPosts(c.Context(), filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar postingan",
			nil,
		))
	}

	meta := &domain.Meta{
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Total:   total,
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(posts, meta))
}

func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	post, err := h.svc.GetPostByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Postingan tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil postingan",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(post, nil))
}

func (h *PostHandler) Create(c *fiber.Ctx) error {
	var userID string
	if user, ok := c.Locals("user").(*domain.User); ok && user != nil {
		userID = user.ID
	} else if userResp, ok := c.Locals("user").(*domain.UserResponse); ok && userResp != nil {
		userID = userResp.ID
	} else {
		return c.Status(http.StatusUnauthorized).JSON(domain.NewErrorResponse(
			"UNAUTHORIZED",
			"Autentikasi diperlukan",
			nil,
		))
	}

	var req domain.PostCreateRequest
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
			"Data postingan tidak valid",
			fields,
		))
	}

	post, err := h.svc.CreatePost(c.Context(), userID, req)
	if err != nil {
		if errors.Is(err, repository.ErrPostConflict) {
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

	return c.Status(http.StatusCreated).JSON(domain.NewSuccessResponse(post, nil))
}

func (h *PostHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.PostUpdateRequest
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
			"Data postingan tidak valid",
			fields,
		))
	}

	post, err := h.svc.UpdatePost(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Postingan tidak ditemukan",
				nil,
			))
		}
		if errors.Is(err, repository.ErrPostConflict) {
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

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(post, nil))
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeletePost(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Postingan tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"DELETE_FAILED",
			"Gagal menghapus postingan",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Postingan berhasil dihapus",
	}, nil))
}
