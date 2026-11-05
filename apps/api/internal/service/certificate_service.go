package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/sanitize"
	"portfolio/api/internal/storage"
)

var (
	ErrInvalidFileMime = errors.New("file_mime harus salah satu dari: application/pdf, image/png, image/jpeg, image/webp")
)

var allowedCertificateMimes = map[string]bool{
	"application/pdf": true,
	"image/png":       true,
	"image/jpeg":      true,
	"image/webp":      true,
}

type CertificateService interface {
	ListCertificates(ctx context.Context, onlyVisible bool) ([]domain.Certificate, error)
	GetCertificateByID(ctx context.Context, id string) (*domain.Certificate, error)
	CreateCertificate(ctx context.Context, req domain.CertificateCreateRequest) (*domain.Certificate, error)
	UpdateCertificate(ctx context.Context, id string, req domain.CertificateUpdateRequest) (*domain.Certificate, error)
	DeleteCertificate(ctx context.Context, id string) error
}

type certificateService struct {
	repo repository.CertificateRepository
}

func NewCertificateService(repo repository.CertificateRepository) CertificateService {
	return &certificateService{repo: repo}
}

func (s *certificateService) ListCertificates(ctx context.Context, onlyVisible bool) ([]domain.Certificate, error) {
	return s.repo.List(ctx, onlyVisible)
}

func (s *certificateService) GetCertificateByID(ctx context.Context, id string) (*domain.Certificate, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *certificateService) CreateCertificate(ctx context.Context, req domain.CertificateCreateRequest) (*domain.Certificate, error) {
	mime := strings.ToLower(strings.TrimSpace(req.FileMime))
	if !allowedCertificateMimes[mime] {
		return nil, ErrInvalidFileMime
	}

	if req.CredentialURL != nil && !sanitize.ValidateHTTPSURL(*req.CredentialURL) {
		return nil, fmt.Errorf("credential_url: %w", ErrInvalidURL)
	}

	if req.ExpiresAt != nil && req.ExpiresAt.Before(req.IssuedAt) {
		return nil, ErrInvalidDateSpan
	}

	isVisible := true
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}

	cert := &domain.Certificate{
		Title:         sanitize.NormalizeText(req.Title),
		Issuer:        sanitize.NormalizeText(req.Issuer),
		CredentialID:  req.CredentialID,
		CredentialURL: req.CredentialURL,
		IssuedAt:      req.IssuedAt,
		ExpiresAt:     req.ExpiresAt,
		FilePath:      storage.SanitizeStoragePath(req.FilePath),
		FileMime:      mime,
		ThumbnailPath: storage.SanitizeStoragePathPtr(req.ThumbnailPath),
		IsVisible:     isVisible,
		SortOrder:     req.SortOrder,
	}

	if err := s.repo.Create(ctx, cert); err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	return cert, nil
}

func (s *certificateService) UpdateCertificate(ctx context.Context, id string, req domain.CertificateUpdateRequest) (*domain.Certificate, error) {
	cert, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mime := strings.ToLower(strings.TrimSpace(req.FileMime))
	if !allowedCertificateMimes[mime] {
		return nil, ErrInvalidFileMime
	}

	if req.CredentialURL != nil && !sanitize.ValidateHTTPSURL(*req.CredentialURL) {
		return nil, fmt.Errorf("credential_url: %w", ErrInvalidURL)
	}

	if req.ExpiresAt != nil && req.ExpiresAt.Before(req.IssuedAt) {
		return nil, ErrInvalidDateSpan
	}

	isVisible := cert.IsVisible
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}

	cert.Title = sanitize.NormalizeText(req.Title)
	cert.Issuer = sanitize.NormalizeText(req.Issuer)
	cert.CredentialID = req.CredentialID
	cert.CredentialURL = req.CredentialURL
	cert.IssuedAt = req.IssuedAt
	cert.ExpiresAt = req.ExpiresAt
	cert.FilePath = storage.SanitizeStoragePath(req.FilePath)
	cert.FileMime = mime
	cert.ThumbnailPath = storage.SanitizeStoragePathPtr(req.ThumbnailPath)
	cert.IsVisible = isVisible
	cert.SortOrder = req.SortOrder

	if err := s.repo.Update(ctx, id, cert); err != nil {
		return nil, err
	}

	return cert, nil
}

func (s *certificateService) DeleteCertificate(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
