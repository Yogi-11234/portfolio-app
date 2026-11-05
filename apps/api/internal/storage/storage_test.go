package storage

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---- helpers ---------------------------------------------------------------

// makeTempProvider creates a LocalStorageProvider backed by a temp directory.
func makeTempProvider(t *testing.T) *LocalStorageProvider {
	t.Helper()
	dir := t.TempDir()
	p, err := NewLocalStorageProvider(dir)
	if err != nil {
		t.Fatalf("NewLocalStorageProvider: %v", err)
	}
	return p
}

// fakeJPEGHeader is the 3-byte JPEG magic: FF D8 FF.
var fakeJPEGHeader = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}

// fakePNGHeader is the 8-byte PNG signature.
var fakePNGHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

// fakePDFHeader starts with %PDF-.
var fakePDFHeader = []byte{0x25, 0x50, 0x44, 0x46, 0x2D}

// makeFileHeader creates a *multipart.FileHeader from raw bytes + filename.
func makeFileHeader(t *testing.T, filename string, body []byte) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		`form-data; name="file"; filename="`+filename+`"`)
	h.Set("Content-Type", "application/octet-stream")

	fw, err := mw.CreatePart(h)
	if err != nil {
		t.Fatalf("CreatePart: %v", err)
	}
	if _, err := fw.Write(body); err != nil {
		t.Fatalf("Write: %v", err)
	}
	mw.Close()

	mr := multipart.NewReader(&buf, mw.Boundary())
	form, err := mr.ReadForm(10 << 20)
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}
	files, ok := form.File["file"]
	if !ok || len(files) == 0 {
		t.Fatal("no file in form")
	}
	return files[0]
}

// ---- tests -----------------------------------------------------------------

func TestSave_ValidJPEG(t *testing.T) {
	p := makeTempProvider(t)
	fh := makeFileHeader(t, "photo.jpg", append(fakeJPEGHeader, make([]byte, 100)...))
	fh.Size = int64(len(fakeJPEGHeader) + 100)

	res, err := p.Save(KindAvatar, fh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(res.RelativePath, "avatars/") {
		t.Fatalf("expected path to start with avatars/, got %q", res.RelativePath)
	}
	if !strings.HasSuffix(res.RelativePath, ".jpg") {
		t.Fatalf("expected .jpg extension, got %q", res.RelativePath)
	}
	// Verify file actually exists on disk.
	if _, err := os.Stat(p.AbsPath(res.RelativePath)); err != nil {
		t.Fatalf("file not found on disk: %v", err)
	}
}

func TestSave_ValidPNG(t *testing.T) {
	p := makeTempProvider(t)
	fh := makeFileHeader(t, "image.png", append(fakePNGHeader, make([]byte, 50)...))
	fh.Size = int64(len(fakePNGHeader) + 50)

	res, err := p.Save(KindProject, fh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MIME != "image/png" {
		t.Fatalf("expected image/png, got %s", res.MIME)
	}
}

func TestSave_ValidPDF(t *testing.T) {
	p := makeTempProvider(t)
	body := append(fakePDFHeader, make([]byte, 200)...)
	fh := makeFileHeader(t, "resume.pdf", body)
	fh.Size = int64(len(body))

	res, err := p.Save(KindResume, fh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(res.RelativePath, ".pdf") {
		t.Fatalf("expected .pdf extension")
	}
}

func TestSave_RejectFileTooLarge(t *testing.T) {
	p := makeTempProvider(t)
	fh := makeFileHeader(t, "big.jpg", fakeJPEGHeader)
	fh.Size = maxUploadBytes + 1 // trigger size guard

	_, err := p.Save(KindAvatar, fh)
	if err != ErrFileTooLarge {
		t.Fatalf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestSave_RejectSVGByExtension(t *testing.T) {
	p := makeTempProvider(t)
	svgBody := []byte(`<svg xmlns="http://www.w3.org/2000/svg"><rect/></svg>`)
	fh := makeFileHeader(t, "icon.svg", svgBody)
	fh.Size = int64(len(svgBody))

	_, err := p.Save(KindProject, fh)
	if err != ErrExtensionNotAllowed {
		t.Fatalf("expected ErrExtensionNotAllowed for .svg, got %v", err)
	}
}

func TestSave_RejectDisallowedExtension(t *testing.T) {
	p := makeTempProvider(t)
	fh := makeFileHeader(t, "script.exe", []byte("MZ\x90\x00"))
	fh.Size = 4

	_, err := p.Save(KindProject, fh)
	if err != ErrExtensionNotAllowed {
		t.Fatalf("expected ErrExtensionNotAllowed, got %v", err)
	}
}

func TestSave_RejectMIMEMismatch(t *testing.T) {
	// File has .jpg extension but contains plain text body — DetectContentType
	// will return "text/plain; charset=utf-8", not "image/jpeg".
	p := makeTempProvider(t)
	fh := makeFileHeader(t, "fake.jpg", []byte("Hello, I am not a JPEG!"))
	fh.Size = int64(len("Hello, I am not a JPEG!"))

	_, err := p.Save(KindAvatar, fh)
	if err != ErrMIMENotAllowed {
		t.Fatalf("expected ErrMIMENotAllowed, got %v", err)
	}
}

func TestSave_FilenameIsUUID(t *testing.T) {
	p := makeTempProvider(t)
	fh := makeFileHeader(t, "../../../../etc/passwd.jpg", append(fakeJPEGHeader, make([]byte, 20)...))
	fh.Size = int64(len(fakeJPEGHeader) + 20)

	res, err := p.Save(KindAvatar, fh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Original filename must NOT appear in the stored path.
	if strings.Contains(res.RelativePath, "passwd") || strings.Contains(res.RelativePath, "..") {
		t.Fatalf("stored path leaks original filename: %q", res.RelativePath)
	}
}

func TestDelete_Idempotent(t *testing.T) {
	p := makeTempProvider(t)
	// Deleting a non-existent file should not error.
	if err := p.Delete("avatars/nonexistent.jpg"); err != nil {
		t.Fatalf("unexpected error deleting non-existent file: %v", err)
	}
}

func TestDelete_ValidFile(t *testing.T) {
	p := makeTempProvider(t)
	// Create a real file and then delete it.
	dir := filepath.Join(p.baseDir, "avatars")
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(dir, "test*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	rel := filepath.ToSlash(filepath.Join("avatars", filepath.Base(f.Name())))
	if err := p.Delete(rel); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(f.Name()); !os.IsNotExist(err) {
		t.Fatal("file should have been deleted")
	}
}

func TestDelete_PathTraversal(t *testing.T) {
	p := makeTempProvider(t)
	err := p.Delete("../../../etc/passwd")
	if err != ErrPathTraversal {
		t.Fatalf("expected ErrPathTraversal, got %v", err)
	}
}

func TestPublicURL(t *testing.T) {
	p := makeTempProvider(t)
	p.baseURL = "https://kpbuhrpgnrqmwrzzxols.supabase.co/storage/v1/object/public/portfolio_assets"
	p.driver = "supabase"
	url := p.PublicURL("avatars/abc.jpg")
	expected := "https://kpbuhrpgnrqmwrzzxols.supabase.co/storage/v1/object/public/portfolio_assets/avatars/abc.jpg"
	if url != expected {
		t.Fatalf("unexpected URL: %s, expected: %s", url, expected)
	}

	// Test local driver
	pLocal, _ := NewLocalStorageProvider(t.TempDir(), "http://localhost:8080", "local")
	localURL := pLocal.PublicURL("avatars/abc.jpg")
	if localURL != "http://localhost:8080/uploads/avatars/abc.jpg" {
		t.Fatalf("unexpected local URL: %s", localURL)
	}

	// Test SanitizeStoragePath
	raw := "https://kpbuhrpgnrqmwrzzxols.supabase.co/storage/v1/object/public/portfolio_assets/projects/pic.png"
	sanitized := SanitizeStoragePath(raw)
	if sanitized != "projects/pic.png" {
		t.Fatalf("unexpected sanitized path: %s", sanitized)
	}

	rawSupabaseDirect := "https://kpbuhrpgnrqmwrzzxols.supabase.co/storage/v1/object/portfolio_assets/projects/pic.png"
	sanitizedDirect := SanitizeStoragePath(rawSupabaseDirect)
	if sanitizedDirect != "projects/pic.png" {
		t.Fatalf("unexpected sanitized direct path: %s", sanitizedDirect)
	}

	rawLocal := "http://localhost:5173/uploads/certificates/cert.pdf"
	sanitizedLocal := SanitizeStoragePath(rawLocal)
	if sanitizedLocal != "certificates/cert.pdf" {
		t.Fatalf("unexpected sanitized local path: %s", sanitizedLocal)
	}
}

func TestSupabaseStorageProvider_SaveAndDelete(t *testing.T) {
	var capturedMethod, capturedPath, capturedAuth, capturedAPIKey, capturedContentType, capturedUpsert string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		capturedAuth = r.Header.Get("Authorization")
		capturedAPIKey = r.Header.Get("apikey")
		capturedContentType = r.Header.Get("Content-Type")
		capturedUpsert = r.Header.Get("x-upsert")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Key":"portfolio_assets/projects/test.png"}`))
	}))
	defer mockServer.Close()

	provider := NewSupabaseStorageProvider(
		mockServer.URL,
		"portfolio_assets",
		"test-anon-key",
		"",
	)
	provider.SetHTTPClient(mockServer.Client())

	fh := makeFileHeader(t, "diagram.png", append(fakePNGHeader, make([]byte, 50)...))
	fh.Size = int64(len(fakePNGHeader) + 50)

	res, err := provider.Save(KindProject, fh)
	if err != nil {
		t.Fatalf("Supabase Save failed: %v", err)
	}

	if !strings.HasPrefix(res.RelativePath, "projects/") {
		t.Fatalf("expected relative path to start with projects/, got: %s", res.RelativePath)
	}
	if capturedMethod != "POST" {
		t.Fatalf("expected POST method, got: %s", capturedMethod)
	}
	if !strings.HasPrefix(capturedPath, "/storage/v1/object/portfolio_assets/projects/") {
		t.Fatalf("expected path to start with /storage/v1/object/portfolio_assets/projects/, got: %s", capturedPath)
	}
	if capturedAuth != "Bearer test-anon-key" {
		t.Fatalf("expected Authorization header Bearer test-anon-key, got: %s", capturedAuth)
	}
	if capturedAPIKey != "test-anon-key" {
		t.Fatalf("expected apikey header test-anon-key, got: %s", capturedAPIKey)
	}
	if capturedContentType != "image/png" {
		t.Fatalf("expected Content-Type image/png, got: %s", capturedContentType)
	}
	if capturedUpsert != "true" {
		t.Fatalf("expected x-upsert true, got: %s", capturedUpsert)
	}

	// Test Delete
	err = provider.Delete(res.RelativePath)
	if err != nil {
		t.Fatalf("Supabase Delete failed: %v", err)
	}
	if capturedMethod != "DELETE" {
		t.Fatalf("expected DELETE method, got: %s", capturedMethod)
	}

	// Test PublicURL
	pubURL := provider.PublicURL(res.RelativePath)
	expectedPrefix := mockServer.URL + "/storage/v1/object/public/portfolio_assets/projects/"
	if !strings.HasPrefix(pubURL, expectedPrefix) {
		t.Fatalf("expected public URL prefix %s, got: %s", expectedPrefix, pubURL)
	}
}

// TestSave_RelativePathForwardSlashes ensures the stored path uses "/" even on Windows.
func TestSave_RelativePathForwardSlashes(t *testing.T) {
	p := makeTempProvider(t)
	fh := makeFileHeader(t, "pic.png", append(fakePNGHeader, make([]byte, 10)...))
	fh.Size = int64(len(fakePNGHeader) + 10)

	res, err := p.Save(KindPost, fh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(res.RelativePath, `\`) {
		t.Fatalf("relative path must use forward slashes, got: %q", res.RelativePath)
	}
}

// Compile-time checks: Providers must satisfy StorageProvider.
var _ StorageProvider = (*LocalStorageProvider)(nil)
var _ StorageProvider = (*SupabaseStorageProvider)(nil)

// Satisfy io import (used indirectly via io.Copy in storage.go).
var _ io.Writer = io.Discard
