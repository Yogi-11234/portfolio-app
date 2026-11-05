package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
)

type mockProjectRepo struct {
	projects []domain.Project
}

func (m *mockProjectRepo) List(ctx context.Context, f domain.ProjectFilter) ([]domain.Project, int, error) {
	return m.projects, len(m.projects), nil
}

func (m *mockProjectRepo) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	for _, p := range m.projects {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, repository.ErrProjectNotFound
}

func (m *mockProjectRepo) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	for _, p := range m.projects {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, repository.ErrProjectNotFound
}

func (m *mockProjectRepo) Create(ctx context.Context, p *domain.Project, skillIDs []string) error {
	p.ID = "proj-uuid-" + p.Slug
	m.projects = append(m.projects, *p)
	return nil
}

func (m *mockProjectRepo) Update(ctx context.Context, id string, p *domain.Project, skillIDs []string) error {
	for i, existing := range m.projects {
		if existing.ID == id {
			m.projects[i] = *p
			return nil
		}
	}
	return repository.ErrProjectNotFound
}

func (m *mockProjectRepo) UpdateSkills(ctx context.Context, projectID string, skillIDs []string) error {
	return nil
}

func (m *mockProjectRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func TestProjectService(t *testing.T) {
	mockProjRepo := &mockProjectRepo{}
	mockCatRepo := &mockCategoryRepo{
		categories: []domain.Category{
			{ID: "cat-proj", Name: "Web", Slug: "web", Type: "project"},
			{ID: "cat-post", Name: "Blog", Slug: "blog", Type: "post"},
		},
	}
	mockSkillRepo := &mockSkillRepo{
		skills: []domain.Skill{
			{ID: "skill-go", Name: "Go", Slug: "go"},
		},
	}

	catSvc := NewCategoryService(mockCatRepo)
	skillSvc := NewSkillService(mockSkillRepo)
	projSvc := NewProjectService(mockProjRepo, catSvc, skillSvc)
	ctx := context.Background()

	// 1. Invalid status
	_, err := projSvc.CreateProject(ctx, domain.ProjectCreateRequest{
		Title:   "My Project",
		Summary: "Summary",
		Status:  "invalid_status",
	})
	if err != ErrInvalidStatus {
		t.Errorf("Expected ErrInvalidStatus, got: %v", err)
	}

	// 2. Insecure repo URL
	insecureURL := "http://github.com/repo"
	_, err = projSvc.CreateProject(ctx, domain.ProjectCreateRequest{
		Title:   "My Project",
		Summary: "Summary",
		Status:  "published",
		RepoURL: &insecureURL,
	})
	if err == nil || !strings.Contains(err.Error(), "https://") {
		t.Errorf("Expected error for non-https repo URL, got: %v", err)
	}

	// 3. Category of wrong type (post category for project)
	postCatID := "cat-post"
	_, err = projSvc.CreateProject(ctx, domain.ProjectCreateRequest{
		Title:      "My Project",
		Summary:    "Summary",
		Status:     "published",
		CategoryID: &postCatID,
	})
	if err == nil || !strings.Contains(err.Error(), "bukan untuk tipe") {
		t.Errorf("Expected category type mismatch error, got: %v", err)
	}

	// 4. Invalid date span
	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err = projSvc.CreateProject(ctx, domain.ProjectCreateRequest{
		Title:     "My Project",
		Summary:   "Summary",
		Status:    "published",
		StartedAt: &start,
		EndedAt:   &end,
	})
	if err != ErrInvalidDateSpan {
		t.Errorf("Expected ErrInvalidDateSpan, got: %v", err)
	}

	// 5. Valid Project creation
	validRepo := "https://github.com/valid/repo"
	projCatID := "cat-proj"
	validProj, err := projSvc.CreateProject(ctx, domain.ProjectCreateRequest{
		Title:         "E-Commerce Platform",
		Summary:       "Scalable shop with microservices",
		DescriptionMD: "Built with **Go Fiber** <script>bad()</script>",
		Status:        "published",
		CategoryID:    &projCatID,
		SkillIDs:      []string{"skill-go"},
		RepoURL:       &validRepo,
	})
	if err != nil {
		t.Fatalf("Expected create project to succeed, got: %v", err)
	}
	if validProj.Slug != "e-commerce-platform" {
		t.Errorf("Expected slug 'e-commerce-platform', got '%s'", validProj.Slug)
	}
	if strings.Contains(validProj.DescriptionHTML, "<script") {
		t.Error("Expected description_html to be sanitized")
	}
}
