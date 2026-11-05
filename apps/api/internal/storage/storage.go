// Package storage defines the StorageProvider contract and a local-disk
// implementation. All validation (extension whitelist, magic-bytes, size,
// SVG rejection, path-traversal) lives here so it applies regardless of
// which handler calls Save().
package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ---- sentinel errors -------------------------------------------------------

var (
	ErrFileTooLarge        = errors.New("ukuran file melebihi batas maksimum 2 MB")
	ErrExtensionNotAllowed = errors.New("tipe file tidak diizinkan")
	ErrMIMENotAllowed      = errors.New("format file tidak didukung (magic bytes tidak valid)")
	ErrSVGNotAllowed       = errors.New("file SVG tidak diizinkan demi keamanan")
	ErrPathTraversal       = errors.New("path file tidak valid")
)

// ---- allowed types ---------------------------------------------------------

// allowedExtensions maps lowercase file extension → canonical MIME type.
// SVG is intentionally absent.
var allowedExtensions = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
	".pdf":  "application/pdf",
}

// allowedMIMEs is the set of MIME types we accept (derived from
// http.DetectContentType magic-bytes sniffing).
var allowedMIMEs = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/webp":      true,
	"application/pdf": true,
}

// UploadKind categorises what is being uploaded so the provider can place
// the file in the right sub-directory and apply any per-kind rules.
type UploadKind string

const (
	KindAvatar      UploadKind = "avatars"
	KindResume      UploadKind = "resumes"
	KindProject     UploadKind = "projects"
	KindCertificate UploadKind = "certificates"
	KindPost        UploadKind = "posts"
	KindSkill       UploadKind = "skills"
)

// ValidKinds is the complete set of accepted upload kinds.
var ValidKinds = map[UploadKind]bool{
	KindAvatar:      true,
	KindResume:      true,
	KindProject:     true,
	KindCertificate: true,
	KindPost:        true,
	KindSkill:       true,
}

// ---- SaveResult ------------------------------------------------------------

// SaveResult contains the metadata of a successfully stored file.
type SaveResult struct {
	// RelativePath is the path relative to StorageLocalDir,
	// e.g. "avatars/3f2504e0-4f89.jpg". This is what gets stored in the DB.
	RelativePath string
	// MIME is the canonical MIME type detected via magic bytes.
	MIME string
	// Size is the file size in bytes.
	Size int64
}

// ---- StorageProvider interface ---------------------------------------------

// StorageProvider is the single abstraction for all file I/O operations.
// Swapping local disk for Appwrite Cloud means only implementing this interface.
type StorageProvider interface {
	// Save validates and stores the uploaded file.
	// kind must be one of the UploadKind constants.
	// Returns a SaveResult on success; a typed sentinel error on failure.
	Save(kind UploadKind, fh *multipart.FileHeader) (*SaveResult, error)

	// Delete removes a file identified by its relative path.
	// Returns nil if the file does not exist (idempotent).
	Delete(relativePath string) error

	// AbsPath converts a relative storage path to an absolute filesystem path.
	AbsPath(relativePath string) string

	// PublicURL returns the URL path for serving the file
	// (e.g. "/uploads/avatars/abc.jpg").
	PublicURL(relativePath string) string
}

// ---- LocalStorageProvider --------------------------------------------------

const maxUploadBytes = 2 * 1024 * 1024 // 2 MB

var (
	globalBaseURL = "https://kpbuhrpgnrqmwrzzxols.supabase.co/storage/v1/object/public/portfolio_assets"
	globalDriver  = "supabase"
)

// SetGlobalStorageConfig sets the global base URL and storage driver used by FormatMediaURL.
func SetGlobalStorageConfig(baseURL, driver string) {
	if baseURL != "" {
		globalBaseURL = strings.TrimRight(baseURL, "/")
	}
	if driver != "" {
		globalDriver = strings.ToLower(strings.TrimSpace(driver))
	}
}

// BuildStorageURL formats a relative storage path into a full public URL based on the driver and base URL.
// 1. If path is empty, returns "".
// 2. If path already begins with "http://" or "https://", returns path untouched.
// 3. If relative path (e.g. "projects/399fd81...png"), strips any leading slashes or "/uploads/" prefix,
//    and joins directly with appBaseURL:
//    fmt.Sprintf("%s/%s", appBaseURL, relativePath) -> produces "https://kpbuhrpgnrqmwrzzxols.supabase.co/storage/v1/object/public/portfolio_assets/projects/399fd81...png"
func BuildStorageURL(appBaseURL, storageDriver, path string) string {
	if path == "" {
		return ""
	}
	clean := strings.TrimSpace(path)
	if strings.HasPrefix(clean, "http://") || strings.HasPrefix(clean, "https://") {
		return clean
	}

	clean = SanitizeStoragePath(clean)

	baseURL := strings.TrimRight(appBaseURL, "/")
	driver := strings.ToLower(strings.TrimSpace(storageDriver))

	if driver == "supabase" || driver == "cloud" {
		return fmt.Sprintf("%s/%s", baseURL, clean)
	}

	// Local driver
	if strings.HasSuffix(baseURL, "/uploads") {
		return fmt.Sprintf("%s/%s", baseURL, clean)
	}
	return fmt.Sprintf("%s/uploads/%s", baseURL, clean)
}

// FormatMediaURL formats a relative path or full URL using the configured global base URL and driver.
func FormatMediaURL(path string) string {
	return BuildStorageURL(globalBaseURL, globalDriver, path)
}

// FormatMediaURLPtr formats a string pointer path into a full public URL.
func FormatMediaURLPtr(path *string) *string {
	if path == nil || *path == "" {
		return nil
	}
	formatted := FormatMediaURL(*path)
	if formatted == "" {
		return nil
	}
	return &formatted
}

// SanitizeStoragePath cleans any full URL (Supabase, localhost, etc.) or "/uploads/" prefix from input,
// ensuring ONLY the relative path (e.g. "projects/nama-file.png", "certificates/nama-file.pdf") is stored in the database.
func SanitizeStoragePath(path string) string {
	if path == "" {
		return ""
	}
	clean := strings.TrimSpace(path)

	// If full URL, strip scheme and domain/base prefix
	if idx := strings.Index(clean, "://"); idx != -1 {
		afterScheme := clean[idx+3:]
		if slashIdx := strings.Index(afterScheme, "/"); slashIdx != -1 {
			clean = afterScheme[slashIdx:]
		}
	}

	// Strip common Supabase storage paths (with public or direct object path)
	if idx := strings.Index(clean, "/storage/v1/object/public/"); idx != -1 {
		rest := clean[idx+len("/storage/v1/object/public/"):]
		if slashIdx := strings.Index(rest, "/"); slashIdx != -1 {
			clean = rest[slashIdx+1:]
		} else {
			clean = rest
		}
	} else if idx := strings.Index(clean, "/storage/v1/object/"); idx != -1 {
		rest := clean[idx+len("/storage/v1/object/"):]
		if slashIdx := strings.Index(rest, "/"); slashIdx != -1 {
			clean = rest[slashIdx+1:]
		} else {
			clean = rest
		}
	}

	clean = strings.TrimPrefix(clean, "/")
	clean = strings.TrimPrefix(clean, "uploads/")
	clean = strings.TrimPrefix(clean, "/")

	return clean
}

// SanitizeStoragePathPtr sanitizes a pointer to a storage path string.
func SanitizeStoragePathPtr(path *string) *string {
	if path == nil || *path == "" {
		return nil
	}
	sanitized := SanitizeStoragePath(*path)
	if sanitized == "" {
		return nil
	}
	return &sanitized
}

// preparedUpload contains validated upload components.
type preparedUpload struct {
	file         multipart.File
	filename     string
	relativePath string
	mime         string
	size         int64
}

func prepareUpload(kind UploadKind, fh *multipart.FileHeader) (*preparedUpload, error) {
	// 1. Validate upload kind.
	if !ValidKinds[kind] {
		return nil, fmt.Errorf("storage: unknown upload kind %q", kind)
	}

	// 2. Size guard (early, before reading the whole file).
	if fh.Size > maxUploadBytes {
		return nil, ErrFileTooLarge
	}

	// 3. Extension whitelist check.
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if _, ok := allowedExtensions[ext]; !ok {
		return nil, ErrExtensionNotAllowed
	}

	// 4. Open the upload for reading.
	src, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("storage: cannot open uploaded file: %w", err)
	}

	// 5. Read first 512 bytes for magic-bytes sniffing (http.DetectContentType).
	sniff := make([]byte, 512)
	n, err := src.Read(sniff)
	if err != nil && !errors.Is(err, io.EOF) {
		_ = src.Close()
		return nil, fmt.Errorf("storage: cannot read file header: %w", err)
	}
	detectedMIME := http.DetectContentType(sniff[:n])

	// Normalise: http.DetectContentType may return "image/jpeg; charset=..." etc.
	if idx := strings.Index(detectedMIME, ";"); idx != -1 {
		detectedMIME = strings.TrimSpace(detectedMIME[:idx])
	}

	// 6. SVG rejection — check both extension and MIME (SVG can masquerade as text/html).
	if ext == ".svg" || strings.Contains(detectedMIME, "svg") {
		_ = src.Close()
		return nil, ErrSVGNotAllowed
	}

	// 7. MIME whitelist (magic bytes).
	if !allowedMIMEs[detectedMIME] {
		_ = src.Close()
		return nil, ErrMIMENotAllowed
	}

	// 8. Seek back to start so we copy the full file.
	if seeker, ok := src.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			_ = src.Close()
			return nil, fmt.Errorf("storage: cannot seek file: %w", err)
		}
	} else {
		// Re-open if the implementation doesn't support seeking.
		_ = src.Close()
		src, err = fh.Open()
		if err != nil {
			return nil, fmt.Errorf("storage: cannot re-open uploaded file: %w", err)
		}
	}

	// 9. Generate UUID-based filename (prevents collisions and enumeration).
	uuidBytes := make([]byte, 16)
	if _, err := rand.Read(uuidBytes); err != nil {
		_ = src.Close()
		return nil, fmt.Errorf("storage: cannot generate file name: %w", err)
	}
	filename := hex.EncodeToString(uuidBytes) + ext
	relativePath := filepath.ToSlash(filepath.Join(string(kind), filename))

	return &preparedUpload{
		file:         src,
		filename:     filename,
		relativePath: relativePath,
		mime:         detectedMIME,
		size:         fh.Size,
	}, nil
}

type LocalStorageProvider struct {
	baseDir string // absolute path to storage root (e.g. "D:\porto project\storage\uploads")
	baseURL string
	driver  string
}

// NewLocalStorageProvider creates a provider rooted at baseDir with optional baseURL and driver.
// It creates the directory if it does not exist.
func NewLocalStorageProvider(baseDir string, opts ...string) (*LocalStorageProvider, error) {
	abs, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("storage: cannot resolve base dir: %w", err)
	}
	if err := os.MkdirAll(abs, 0750); err != nil {
		return nil, fmt.Errorf("storage: cannot create base dir %q: %w", abs, err)
	}
	baseURL := globalBaseURL
	driver := globalDriver
	if len(opts) > 0 && opts[0] != "" {
		baseURL = opts[0]
	}
	if len(opts) > 1 && opts[1] != "" {
		driver = opts[1]
	}
	return &LocalStorageProvider{
		baseDir: abs,
		baseURL: baseURL,
		driver:  driver,
	}, nil
}

// Save validates the file and writes it to baseDir/kind/<uuid>.<ext>.
func (p *LocalStorageProvider) Save(kind UploadKind, fh *multipart.FileHeader) (*SaveResult, error) {
	prep, err := prepareUpload(kind, fh)
	if err != nil {
		return nil, err
	}
	defer prep.file.Close()

	// Build destination path and verify it stays within baseDir (anti path-traversal).
	kindDir := filepath.Join(p.baseDir, string(kind))
	if err := os.MkdirAll(kindDir, 0750); err != nil {
		return nil, fmt.Errorf("storage: cannot create kind dir: %w", err)
	}

	destPath := filepath.Join(kindDir, prep.filename)

	// filepath.Rel check: destPath must be inside kindDir.
	rel, err := filepath.Rel(p.baseDir, destPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return nil, ErrPathTraversal
	}

	// Write file to disk.
	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0640)
	if err != nil {
		return nil, fmt.Errorf("storage: cannot create destination file: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, prep.file)
	if err != nil {
		// Clean up partial write.
		_ = os.Remove(destPath)
		return nil, fmt.Errorf("storage: error writing file: %w", err)
	}

	return &SaveResult{
		RelativePath: prep.relativePath,
		MIME:         prep.mime,
		Size:         written,
	}, nil
}

// Delete removes the file at baseDir/relativePath.
// Returns nil if the file does not exist (idempotent).
func (p *LocalStorageProvider) Delete(relativePath string) error {
	if relativePath == "" {
		return nil
	}

	// Anti path-traversal: relative path must resolve inside baseDir.
	destPath := filepath.Join(p.baseDir, filepath.FromSlash(relativePath))
	rel, err := filepath.Rel(p.baseDir, destPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return ErrPathTraversal
	}

	err = os.Remove(destPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// AbsPath converts a relative storage path to an absolute filesystem path.
func (p *LocalStorageProvider) AbsPath(relativePath string) string {
	return filepath.Join(p.baseDir, filepath.FromSlash(relativePath))
}

// PublicURL returns the public URL for serving the file according to the storage driver and baseURL.
func (p *LocalStorageProvider) PublicURL(relativePath string) string {
	baseURL := p.baseURL
	if baseURL == "" {
		baseURL = globalBaseURL
	}
	driver := p.driver
	if driver == "" {
		driver = globalDriver
	}
	return BuildStorageURL(baseURL, driver, relativePath)
}
