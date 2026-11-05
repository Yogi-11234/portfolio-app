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

type CertificateHandler struct {
	svc service.CertificateService
}

func NewCertificateHandler(svc service.CertificateService) *CertificateHandler {
	return &CertificateHandler{svc: svc}
}

func (h *CertificateHandler) ListPublic(c *fiber.Ctx) error {
	certificates, err := h.svc.ListCertificates(c.Context(), true)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar sertifikat",
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(certificates, nil))
}

func (h *CertificateHandler) ListAdmin(c *fiber.Ctx) error {
	certificates, err := h.svc.ListCertificates(c.Context(), false)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil daftar sertifikat",
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(certificates, nil))
}

func (h *CertificateHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	certificate, err := h.svc.GetCertificateByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCertificateNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Sertifikat tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"Gagal mengambil sertifikat",
			nil,
		))
	}
	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(certificate, nil))
}

func (h *CertificateHandler) Create(c *fiber.Ctx) error {
	var req domain.CertificateCreateRequest
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
			"Data sertifikat tidak valid",
			fields,
		))
	}

	cert, err := h.svc.CreateCertificate(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"CREATION_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusCreated).JSON(domain.NewSuccessResponse(cert, nil))
}

func (h *CertificateHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.CertificateUpdateRequest
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
			"Data sertifikat tidak valid",
			fields,
		))
	}

	cert, err := h.svc.UpdateCertificate(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrCertificateNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Sertifikat tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"UPDATE_FAILED",
			err.Error(),
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(cert, nil))
}

func (h *CertificateHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteCertificate(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrCertificateNotFound) {
			return c.Status(http.StatusNotFound).JSON(domain.NewErrorResponse(
				"NOT_FOUND",
				"Sertifikat tidak ditemukan",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"DELETE_FAILED",
			"Gagal menghapus sertifikat",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "Sertifikat berhasil dihapus",
	}, nil))
}
