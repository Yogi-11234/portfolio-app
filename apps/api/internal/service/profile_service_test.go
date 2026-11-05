package service

import (
	"context"
	"strings"
	"testing"

	"portfolio/api/internal/domain"
)

type mockProfileRepo struct {
	profile *domain.Profile
}

func (m *mockProfileRepo) Get(ctx context.Context) (*domain.Profile, error) {
	return m.profile, nil
}

func (m *mockProfileRepo) GetByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	return m.profile, nil
}

func (m *mockProfileRepo) Update(ctx context.Context, userID string, p *domain.Profile) error {
	m.profile = p
	return nil
}

func TestProfileService(t *testing.T) {
	mockRepo := &mockProfileRepo{
		profile: &domain.Profile{
			UserID:   "admin-uuid",
			FullName: "Old Name",
		},
	}
	svc := NewProfileService(mockRepo)
	ctx := context.Background()

	// 1. Insecure URL rejection
	insecureURL := "http://insecure.com"
	_, err := svc.UpdateProfile(ctx, "admin-uuid", domain.ProfileUpdateRequest{
		FullName:  "John Doe",
		Headline:  "Engineer",
		GithubURL: &insecureURL,
	})
	if err == nil || !strings.Contains(err.Error(), "https://") {
		t.Fatalf("Expected error for non-https URL, got: %v", err)
	}

	// 2. Successful update with markdown sanitization
	secureURL := "https://github.com/developer"
	res, err := svc.UpdateProfile(ctx, "admin-uuid", domain.ProfileUpdateRequest{
		FullName:  "John Doe",
		Headline:  "Senior Engineer",
		BioMD:     "Hello **world** with <script>alert(1)</script>",
		GithubURL: &secureURL,
	})
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}
	if res.FullName != "John Doe" {
		t.Errorf("Expected FullName 'John Doe', got '%s'", res.FullName)
	}
	if strings.Contains(res.BioHTML, "<script") {
		t.Errorf("Expected BioHTML to be sanitized, got: %s", res.BioHTML)
	}
}
