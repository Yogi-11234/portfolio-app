package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/storage"
)

// UploadHandler handles all multipart file upload endpoints.
// It relies on the StorageProvider for validation and persistence; handlers
// only deal with HTTP concerns (parsing, response shaping).
type UploadHandler struct {
	store storage.StorageProvider
}

func NewUploadHandler(store storage.StorageProvider) *UploadHandler {
	return &UploadHandler{store: store}
}

// ---- shared internals ------------------------------------------------------

// saveFile is the common path: parse multipart → validate via store.Save →
// return the public URL and relative path.
func (h *UploadHandler) saveFile(c *fiber.Ctx, kind storage.UploadKind) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"NO_FILE",
			"Field 'file' tidak ditemukan dalam request",
			nil,
		))
	}

	result, err := h.store.Save(kind, fh)
	if err != nil {
		slog.Error("Upload file failed", "kind", kind, "error", err)
		return uploadError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(domain.NewSuccessResponse(fiber.Map{
		"path":       result.RelativePath,
		"public_url": h.store.PublicURL(result.RelativePath),
		"mime":       result.MIME,
		"size":       result.Size,
	}, nil))
}

// uploadError maps storage sentinel errors to appropriate HTTP responses.
func uploadError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, storage.ErrFileTooLarge):
		return c.Status(http.StatusRequestEntityTooLarge).JSON(domain.NewErrorResponse(
			"FILE_TOO_LARGE",
			err.Error(),
			nil,
		))
	case errors.Is(err, storage.ErrExtensionNotAllowed),
		errors.Is(err, storage.ErrMIMENotAllowed),
		errors.Is(err, storage.ErrSVGNotAllowed):
		return c.Status(http.StatusUnsupportedMediaType).JSON(domain.NewErrorResponse(
			"UNSUPPORTED_FILE_TYPE",
			err.Error(),
			nil,
		))
	case errors.Is(err, storage.ErrPathTraversal):
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"INVALID_FILE_PATH",
			"Path file tidak valid",
			nil,
		))
	default:
		msg := "Gagal menyimpan file. Coba lagi."
		if err != nil && (strings.Contains(err.Error(), "SUPABASE_KEY") || strings.Contains(err.Error(), "supabase")) {
			msg = err.Error()
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"UPLOAD_FAILED",
			msg,
			nil,
		))
	}
}

// ---- endpoints -------------------------------------------------------------

// UploadAvatar handles POST /api/v1/admin/uploads/avatar
// Accepts: .jpg .jpeg .png .webp (max 2 MB)
func (h *UploadHandler) UploadAvatar(c *fiber.Ctx) error {
	return h.saveFile(c, storage.KindAvatar)
}

// UploadResume handles POST /api/v1/admin/uploads/resume
// Accepts: .pdf only (max 2 MB)
func (h *UploadHandler) UploadResume(c *fiber.Ctx) error {
	return h.saveFile(c, storage.KindResume)
}

// UploadProjectThumbnail handles POST /api/v1/admin/uploads/project
// Accepts: .jpg .jpeg .png .webp (max 2 MB)
func (h *UploadHandler) UploadProjectThumbnail(c *fiber.Ctx) error {
	return h.saveFile(c, storage.KindProject)
}

// UploadCertificateFile handles POST /api/v1/admin/uploads/certificate
// Accepts: .jpg .jpeg .png .webp .pdf (max 2 MB)
func (h *UploadHandler) UploadCertificateFile(c *fiber.Ctx) error {
	return h.saveFile(c, storage.KindCertificate)
}

// UploadPostCover handles POST /api/v1/admin/uploads/post
// Accepts: .jpg .jpeg .png .webp (max 2 MB)
func (h *UploadHandler) UploadPostCover(c *fiber.Ctx) error {
	return h.saveFile(c, storage.KindPost)
}

// UploadSkillIcon handles POST /api/v1/admin/uploads/skill
// Accepts: .jpg .jpeg .png .webp (max 2 MB)
func (h *UploadHandler) UploadSkillIcon(c *fiber.Ctx) error {
	return h.saveFile(c, storage.KindSkill)
}

// DeleteFile handles DELETE /api/v1/admin/uploads
// Body: { "path": "avatars/abc123.jpg" }
// Used when an admin removes an image without deleting the parent resource.
func (h *UploadHandler) DeleteFile(c *fiber.Ctx) error {
	var body struct {
		Path string `json:"path"`
	}
	if err := c.BodyParser(&body); err != nil || body.Path == "" {
		return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
			"INVALID_REQUEST",
			"Field 'path' wajib diisi",
			nil,
		))
	}

	if err := h.store.Delete(body.Path); err != nil {
		if errors.Is(err, storage.ErrPathTraversal) {
			return c.Status(http.StatusBadRequest).JSON(domain.NewErrorResponse(
				"INVALID_FILE_PATH",
				"Path file tidak valid",
				nil,
			))
		}
		return c.Status(http.StatusInternalServerError).JSON(domain.NewErrorResponse(
			"DELETE_FAILED",
			"Gagal menghapus file",
			nil,
		))
	}

	return c.Status(http.StatusOK).JSON(domain.NewSuccessResponse(fiber.Map{
		"message": "File berhasil dihapus",
	}, nil))
}
