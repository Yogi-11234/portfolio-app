package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
)

type mockCertificateRepo struct {
	list []domain.Certificate
}

func (m *mockCertificateRepo) List(ctx context.Context, onlyVisible bool) ([]domain.Certificate, error) {
	var res []domain.Certificate
	for _, c := range m.list {
		if !onlyVisible || c.IsVisible {
			res = append(res, c)
		}
	}
	return res, nil
}

func (m *mockCertificateRepo) GetByID(ctx context.Context, id string) (*domain.Certificate, error) {
	for _, c := range m.list {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, repository.ErrCertificateNotFound
}

func (m *mockCertificateRepo) Create(ctx context.Context, c *domain.Certificate) error {
	c.ID = "cert-uuid-1"
	m.list = append(m.list, *c)
	return nil
}

func (m *mockCertificateRepo) Update(ctx context.Context, id string, c *domain.Certificate) error {
	for i, existing := range m.list {
		if existing.ID == id {
			m.list[i] = *c
			return nil
		}
	}
	return repository.ErrCertificateNotFound
}

func (m *mockCertificateRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func TestCertificateService(t *testing.T) {
	mockRepo := &mockCertificateRepo{}
	svc := NewCertificateService(mockRepo)
	ctx := context.Background()

	// 1. Invalid MIME type (e.g. text/plain or svg)
	_, err := svc.CreateCertificate(ctx, domain.CertificateCreateRequest{
		Title:    "Cloud Architect",
		Issuer:   "Google",
		FilePath: "certificates/test.svg",
		FileMime: "image/svg+xml",
		IssuedAt: time.Now(),
	})
	if err != ErrInvalidFileMime {
		t.Errorf("Expected ErrInvalidFileMime, got: %v", err)
	}

	// 2. Insecure credential URL
	insecureURL := "http://credential.net/123"
	_, err = svc.CreateCertificate(ctx, domain.CertificateCreateRequest{
		Title:         "Cloud Architect",
		Issuer:        "Google",
		FilePath:      "certificates/test.pdf",
		FileMime:      "application/pdf",
		IssuedAt:      time.Now(),
		CredentialURL: &insecureURL,
	})
	if err == nil || !strings.Contains(err.Error(), "https://") {
		t.Errorf("Expected error for non-https credential URL, got: %v", err)
	}

	// 3. Valid Certificate
	validURL := "https://credential.net/123"
	cert, err := svc.CreateCertificate(ctx, domain.CertificateCreateRequest{
		Title:         "Google Cloud Professional Architect",
		Issuer:        "Google Cloud",
		FilePath:      "certificates/gcp-pca.pdf",
		FileMime:      "application/pdf",
		IssuedAt:      time.Now(),
		CredentialURL: &validURL,
	})
	if err != nil {
		t.Fatalf("Expected create certificate to succeed, got: %v", err)
	}
	if cert.FileMime != "application/pdf" {
		t.Errorf("Unexpected mime: %s", cert.FileMime)
	}
}
