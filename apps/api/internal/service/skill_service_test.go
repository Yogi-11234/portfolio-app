package service

import (
	"context"
	"testing"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
)

type mockSkillRepo struct {
	skills []domain.Skill
}

func (m *mockSkillRepo) List(ctx context.Context, onlyVisible bool) ([]domain.Skill, error) {
	var res []domain.Skill
	for _, s := range m.skills {
		if !onlyVisible || s.IsVisible {
			res = append(res, s)
		}
	}
	return res, nil
}

func (m *mockSkillRepo) GetByID(ctx context.Context, id string) (*domain.Skill, error) {
	for _, s := range m.skills {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, repository.ErrSkillNotFound
}

func (m *mockSkillRepo) GetBySlug(ctx context.Context, slug string) (*domain.Skill, error) {
	for _, s := range m.skills {
		if s.Slug == slug {
			return &s, nil
		}
	}
	return nil, repository.ErrSkillNotFound
}

func (m *mockSkillRepo) GetByName(ctx context.Context, name string) (*domain.Skill, error) {
	for _, s := range m.skills {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, repository.ErrSkillNotFound
}

func (m *mockSkillRepo) GetByIDs(ctx context.Context, ids []string) ([]domain.Skill, error) {
	var res []domain.Skill
	for _, id := range ids {
		for _, s := range m.skills {
			if s.ID == id {
				res = append(res, s)
			}
		}
	}
	return res, nil
}

func (m *mockSkillRepo) Create(ctx context.Context, s *domain.Skill) error {
	s.ID = "skill-uuid-" + s.Slug
	m.skills = append(m.skills, *s)
	return nil
}

func (m *mockSkillRepo) Update(ctx context.Context, id string, s *domain.Skill) error {
	for i, existing := range m.skills {
		if existing.ID == id {
			m.skills[i] = *s
			return nil
		}
	}
	return repository.ErrSkillNotFound
}

func (m *mockSkillRepo) Delete(ctx context.Context, id string) error {
	for i, existing := range m.skills {
		if existing.ID == id {
			m.skills = append(m.skills[:i], m.skills[i+1:]...)
			return nil
		}
	}
	return repository.ErrSkillNotFound
}

func TestSkillService(t *testing.T) {
	mockRepo := &mockSkillRepo{
		skills: []domain.Skill{
			{ID: "skill-1", Name: "Go", Slug: "go", GroupName: "language", IsVisible: true},
		},
	}
	svc := NewSkillService(mockRepo)
	ctx := context.Background()

	// 1. Invalid group
	_, err := svc.CreateSkill(ctx, domain.SkillCreateRequest{
		Name:      "Python",
		GroupName: "invalid_group",
	})
	if err != ErrInvalidGroupName {
		t.Errorf("Expected ErrInvalidGroupName, got: %v", err)
	}

	// 2. Invalid Hex Color
	badHex := "blue"
	_, err = svc.CreateSkill(ctx, domain.SkillCreateRequest{
		Name:       "Python",
		GroupName:  "language",
		BrandColor: &badHex,
	})
	if err != ErrInvalidHexColor {
		t.Errorf("Expected ErrInvalidHexColor, got: %v", err)
	}

	// 3. Invalid Level (> 5)
	badLevel := int16(6)
	_, err = svc.CreateSkill(ctx, domain.SkillCreateRequest{
		Name:      "Python",
		GroupName: "language",
		Level:     &badLevel,
	})
	if err != ErrInvalidLevel {
		t.Errorf("Expected ErrInvalidLevel, got: %v", err)
	}

	// 4. Valid Skill
	validHex := "#3776AB"
	validLevel := int16(4)
	newSkill, err := svc.CreateSkill(ctx, domain.SkillCreateRequest{
		Name:       "Python",
		GroupName:  "language",
		BrandColor: &validHex,
		Level:      &validLevel,
	})
	if err != nil {
		t.Fatalf("Expected create skill to succeed, got: %v", err)
	}
	if newSkill.Slug != "python" || newSkill.GroupName != "language" {
		t.Errorf("Unexpected created skill: %+v", newSkill)
	}
}
