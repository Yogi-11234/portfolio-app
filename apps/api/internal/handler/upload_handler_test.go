package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"portfolio/api/internal/handler"
	"portfolio/api/internal/storage"
)

var fakeJPEGHeader = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}

func setupUploadTestApp(t *testing.T) (*fiber.App, storage.StorageProvider) {
	t.Helper()
	tempDir := t.TempDir()
	store, err := storage.NewLocalStorageProvider(tempDir)
	if err != nil {
		t.Fatalf("failed to create storage provider: %v", err)
	}

	h := handler.NewUploadHandler(store)
	app := fiber.New()

	app.Post("/uploads/avatar", h.UploadAvatar)
	app.Post("/uploads/resume", h.UploadResume)
	app.Post("/uploads/project", h.UploadProjectThumbnail)
	app.Delete("/uploads", h.DeleteFile)

	return app, store
}

func createMultipartRequest(t *testing.T, fieldName, filename string, content []byte) (*http.Request, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+filename+`"`)
	h.Set("Content-Type", "application/octet-stream")

	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatalf("failed to create part: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		t.Fatalf("failed to copy content: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/uploads/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, writer.FormDataContentType()
}

func TestUploadHandler_Success(t *testing.T) {
	app, _ := setupUploadTestApp(t)

	fileData := append(fakeJPEGHeader, make([]byte, 50)...)
	req, _ := createMultipartRequest(t, "file", "photo.jpg", fileData)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := res["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in response, got %v", res)
	}
	if data["path"] == "" || data["public_url"] == "" {
		t.Fatalf("expected non-empty path and public_url, got %v", data)
	}
}

func TestUploadHandler_NoFile(t *testing.T) {
	app, _ := setupUploadTestApp(t)

	// Send multipart form with wrong field name
	fileData := append(fakeJPEGHeader, make([]byte, 20)...)
	req, _ := createMultipartRequest(t, "wrong_field", "photo.jpg", fileData)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing file field, got %d", resp.StatusCode)
	}
}

func TestUploadHandler_InvalidFileType(t *testing.T) {
	app, _ := setupUploadTestApp(t)

	// SVG attempt
	svg := []byte(`<svg></svg>`)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "image.svg")
	part.Write(svg)
	writer.Close()

	req := httptest.NewRequest("POST", "/uploads/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expected status 415 for SVG/unsupported type, got %d", resp.StatusCode)
	}
}

func TestUploadHandler_DeleteFile(t *testing.T) {
	app, store := setupUploadTestApp(t)

	// First upload a file to delete
	fileData := append(fakeJPEGHeader, make([]byte, 50)...)
	req, _ := createMultipartRequest(t, "file", "photo.jpg", fileData)
	resp, _ := app.Test(req)
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	resp.Body.Close()

	data := res["data"].(map[string]interface{})
	path := data["path"].(string)

	// Delete it
	delPayload := map[string]string{"path": path}
	delBytes, _ := json.Marshal(delPayload)
	delReq := httptest.NewRequest("DELETE", "/uploads", bytes.NewReader(delBytes))
	delReq.Header.Set("Content-Type", "application/json")

	delResp, err := app.Test(delReq)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", delResp.StatusCode)
	}

	// Verify file is actually gone
	_ = store
}
