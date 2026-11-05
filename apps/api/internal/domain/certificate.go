package domain

import (
	"encoding/json"
	"time"

	"portfolio/api/internal/storage"
)

type Certificate struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Issuer        string     `json:"issuer"`
	CredentialID  *string    `json:"credential_id"`
	CredentialURL *string    `json:"credential_url"`
	IssuedAt      time.Time  `json:"issued_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	FilePath      string     `json:"file_path"`
	FileMime      string     `json:"file_mime"`
	ThumbnailPath *string    `json:"thumbnail_path"`
	IsVisible     bool       `json:"is_visible"`
	SortOrder     int        `json:"sort_order"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (c Certificate) MarshalJSON() ([]byte, error) {
	type Alias Certificate
	filePath := c.FilePath
	if filePath != "" {
		filePath = storage.FormatMediaURL(filePath)
	}
	return json.Marshal(&struct {
		Alias
		Name          string  `json:"name"`
		OrderIndex    int     `json:"order_index"`
		FilePath      string  `json:"file_path"`
		ThumbnailPath *string `json:"thumbnail_path,omitempty"`
	}{
		Alias:         Alias(c),
		Name:          c.Title,
		OrderIndex:    c.SortOrder,
		FilePath:      filePath,
		ThumbnailPath: storage.FormatMediaURLPtr(c.ThumbnailPath),
	})
}

type CertificateCreateRequest struct {
	Title         string     `json:"title" validate:"required,min=2,max=200"`
	Issuer        string     `json:"issuer" validate:"required,min=2,max=160"`
	CredentialID  *string    `json:"credential_id" validate:"omitempty,max=120"`
	CredentialURL *string    `json:"credential_url" validate:"omitempty,max=255"`
	IssuedAt      time.Time  `json:"issued_at" validate:"required"`
	ExpiresAt     *time.Time `json:"expires_at"`
	FilePath      string     `json:"file_path" validate:"required,max=255"`
	FileMime      string     `json:"file_mime" validate:"required,oneof=application/pdf image/png image/jpeg image/webp"`
	ThumbnailPath *string    `json:"thumbnail_path" validate:"omitempty,max=255"`
	IsVisible     *bool      `json:"is_visible"`
	SortOrder     int        `json:"sort_order"`
}

func (r *CertificateCreateRequest) UnmarshalJSON(data []byte) error {
	type Alias CertificateCreateRequest
	var raw struct {
		Alias
		Name       string `json:"name"`
		OrderIndex *int   `json:"order_index"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = CertificateCreateRequest(raw.Alias)
	if r.Title == "" && raw.Name != "" {
		r.Title = raw.Name
	}
	if r.SortOrder == 0 && raw.OrderIndex != nil {
		r.SortOrder = *raw.OrderIndex
	}
	return nil
}

type CertificateUpdateRequest struct {
	Title         string     `json:"title" validate:"required,min=2,max=200"`
	Issuer        string     `json:"issuer" validate:"required,min=2,max=160"`
	CredentialID  *string    `json:"credential_id" validate:"omitempty,max=120"`
	CredentialURL *string    `json:"credential_url" validate:"omitempty,max=255"`
	IssuedAt      time.Time  `json:"issued_at" validate:"required"`
	ExpiresAt     *time.Time `json:"expires_at"`
	FilePath      string     `json:"file_path" validate:"required,max=255"`
	FileMime      string     `json:"file_mime" validate:"required,oneof=application/pdf image/png image/jpeg image/webp"`
	ThumbnailPath *string    `json:"thumbnail_path" validate:"omitempty,max=255"`
	IsVisible     *bool      `json:"is_visible"`
	SortOrder     int        `json:"sort_order"`
}

func (r *CertificateUpdateRequest) UnmarshalJSON(data []byte) error {
	type Alias CertificateUpdateRequest
	var raw struct {
		Alias
		Name       string `json:"name"`
		OrderIndex *int   `json:"order_index"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = CertificateUpdateRequest(raw.Alias)
	if r.Title == "" && raw.Name != "" {
		r.Title = raw.Name
	}
	if r.SortOrder == 0 && raw.OrderIndex != nil {
		r.SortOrder = *raw.OrderIndex
	}
	return nil
}
