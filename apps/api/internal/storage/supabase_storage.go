package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// SupabaseStorageProvider uploads and manages files via Supabase Storage REST API.
type SupabaseStorageProvider struct {
	supabaseURL string
	bucket      string
	key         string
	baseURL     string
	httpClient  *http.Client
}

// NewSupabaseStorageProvider creates a new provider targeting Supabase Storage REST API.
func NewSupabaseStorageProvider(supabaseURL, bucket, key, baseURL string) *SupabaseStorageProvider {
	cleanURL := strings.TrimRight(strings.TrimSpace(supabaseURL), "/")
	if cleanURL == "" {
		cleanURL = "https://kpbuhrpgnrqmwrzzxols.supabase.co"
	}
	cleanBucket := strings.TrimSpace(bucket)
	if cleanBucket == "" {
		cleanBucket = "portfolio_assets"
	}
	cleanKey := strings.TrimSpace(key)
	cleanBaseURL := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if cleanBaseURL == "" {
		cleanBaseURL = fmt.Sprintf("%s/storage/v1/object/public/%s", cleanURL, cleanBucket)
	}

	return &SupabaseStorageProvider{
		supabaseURL: cleanURL,
		bucket:      cleanBucket,
		key:         cleanKey,
		baseURL:     cleanBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetHTTPClient allows overriding the HTTP client (useful for testing).
func (p *SupabaseStorageProvider) SetHTTPClient(client *http.Client) {
	if client != nil {
		p.httpClient = client
	}
}

// Save validates the file and uploads it directly to Supabase Storage REST API:
// POST https://{supabaseURL}/storage/v1/object/{bucket}/{folder}/{filename}
// Returns SaveResult with RelativePath (e.g. "projects/uuid.jpg").
func (p *SupabaseStorageProvider) Save(kind UploadKind, fh *multipart.FileHeader) (*SaveResult, error) {
	if p.key == "" {
		return nil, fmt.Errorf("storage: SUPABASE_KEY belum diisi di file .env (diperlukan untuk header Authorization)")
	}

	prep, err := prepareUpload(kind, fh)
	if err != nil {
		return nil, err
	}
	defer prep.file.Close()

	// Endpoint: {supabaseURL}/storage/v1/object/{bucket}/{folder}/{filename}
	// prep.relativePath is "{kind}/{filename}" e.g. "projects/123456...png"
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", p.supabaseURL, p.bucket, prep.relativePath)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, uploadURL, prep.file)
	if err != nil {
		return nil, fmt.Errorf("storage: failed to build supabase request: %w", err)
	}

	req.Header.Set("Content-Type", prep.mime)
	req.Header.Set("x-upsert", "true")
	if p.key != "" {
		req.Header.Set("Authorization", "Bearer "+p.key)
		req.Header.Set("apikey", p.key)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("storage: failed to upload to supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("storage: supabase upload failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return &SaveResult{
		RelativePath: prep.relativePath,
		MIME:         prep.mime,
		Size:         prep.size,
	}, nil
}

// Delete removes a file from Supabase Storage via REST API.
func (p *SupabaseStorageProvider) Delete(relativePath string) error {
	cleanPath := SanitizeStoragePath(relativePath)
	if cleanPath == "" {
		return nil
	}

	deleteURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", p.supabaseURL, p.bucket, cleanPath)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("storage: failed to build delete request: %w", err)
	}

	if p.key != "" {
		req.Header.Set("Authorization", "Bearer "+p.key)
		req.Header.Set("apikey", p.key)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("storage: failed to delete from supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("storage: supabase delete failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

// AbsPath returns the public URL for Supabase storage as there is no local filesystem path.
func (p *SupabaseStorageProvider) AbsPath(relativePath string) string {
	return p.PublicURL(relativePath)
}

// PublicURL returns the full public URL for the file.
func (p *SupabaseStorageProvider) PublicURL(relativePath string) string {
	baseURL := p.baseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("%s/storage/v1/object/public/%s", p.supabaseURL, p.bucket)
	}
	return BuildStorageURL(baseURL, "supabase", relativePath)
}
