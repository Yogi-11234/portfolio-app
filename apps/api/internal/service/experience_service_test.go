package service

import (
	"context"
	"testing"
	"time"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
)

type mockExperienceRepo struct {
	list []domain.Experience
}

func (m *mockExperienceRepo) List(ctx context.Context, onlyVisible bool, expType string) ([]domain.Experience, error) {
	var res []domain.Experience
	for _, e := range m.list {
		if (!onlyVisible || e.IsVisible) && (expType == "" || e.Type == expType) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (m *mockExperienceRepo) GetByID(ctx context.Context, id string) (*domain.Experience, error) {
	for _, e := range m.list {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, repository.ErrExperienceNotFound
}

func (m *mockExperienceRepo) Create(ctx context.Context, e *domain.Experience) error {
	e.ID = "exp-uuid-1"
	m.list = append(m.list, *e)
	return nil
}

func (m *mockExperienceRepo) Update(ctx context.Context, id string, e *domain.Experience) error {
	for i, existing := range m.list {
		if existing.ID == id {
			m.list[i] = *e
			return nil
		}
	}
	return repository.ErrExperienceNotFound
}

func (m *mockExperienceRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func TestExperienceService(t *testing.T) {
	mockRepo := &mockExperienceRepo{}
	svc := NewExperienceService(mockRepo)
	ctx := context.Background()

	// 1. Invalid type
	_, err := svc.CreateExperience(ctx, domain.ExperienceCreateRequest{
		Type:         "invalid_type",
		Organization: "Tech Corp",
		Role:         "Backend Dev",
		StartDate:    time.Now(),
	})
	if err != ErrInvalidExperienceType {
		t.Errorf("Expected ErrInvalidExperienceType, got: %v", err)
	}

	// 2. Invalid date span
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err = svc.CreateExperience(ctx, domain.ExperienceCreateRequest{
		Type:         "work",
		Organization: "Tech Corp",
		Role:         "Backend Dev",
		StartDate:    start,
		EndDate:      &end,
	})
	if err != ErrInvalidDateSpan {
		t.Errorf("Expected ErrInvalidDateSpan, got: %v", err)
	}

	// 3. Valid experience
	validEnd := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	exp, err := svc.CreateExperience(ctx, domain.ExperienceCreateRequest{
		Type:          "work",
		Organization:  "Tech Corp",
		Role:          "Senior Backend Engineer",
		DescriptionMD: "Led **distributed systems** team",
		StartDate:     start,
		EndDate:       &validEnd,
		IsCurrent:     false,
	})
	if err != nil {
		t.Fatalf("Expected create experience to succeed, got: %v", err)
	}
	if exp.Role != "Senior Backend Engineer" {
		t.Errorf("Unexpected role: %s", exp.Role)
	}
}
