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
	ErrInvalidStatus   = errors.New("status harus berupa draft, published, atau archived")
	ErrInvalidDateSpan = errors.New("ended_at tidak boleh sebelum started_at")
)

type ProjectService interface {
	ListProjects(ctx context.Context, f domain.ProjectFilter) ([]domain.Project, int, error)
	GetProjectByID(ctx context.Context, id string) (*domain.Project, error)
	GetProjectBySlug(ctx context.Context, slug string) (*domain.Project, error)
	CreateProject(ctx context.Context, req domain.ProjectCreateRequest) (*domain.Project, error)
	UpdateProject(ctx context.Context, id string, req domain.ProjectUpdateRequest) (*domain.Project, error)
	UpdateProjectSkills(ctx context.Context, id string, skillIDs []string) error
	DeleteProject(ctx context.Context, id string) error
}

type projectService struct {
	repo        repository.ProjectRepository
	categorySvc CategoryService
	skillSvc    SkillService
}

func NewProjectService(repo repository.ProjectRepository, categorySvc CategoryService, skillSvc SkillService) ProjectService {
	return &projectService{
		repo:        repo,
		categorySvc: categorySvc,
		skillSvc:    skillSvc,
	}
}

func (s *projectService) ListProjects(ctx context.Context, f domain.ProjectFilter) ([]domain.Project, int, error) {
	return s.repo.List(ctx, f)
}

func (s *projectService) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *projectService) GetProjectBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	cleanSlug := strings.ToLower(strings.TrimSpace(slug))
	return s.repo.GetBySlug(ctx, cleanSlug)
}

func (s *projectService) CreateProject(ctx context.Context, req domain.ProjectCreateRequest) (*domain.Project, error) {
	// 1. Status validation
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status != "draft" && status != "published" && status != "archived" {
		return nil, ErrInvalidStatus
	}

	// 2. URLs validation
	if req.RepoURL != nil && !sanitize.ValidateHTTPSURL(*req.RepoURL) {
		return nil, fmt.Errorf("repo_url: %w", ErrInvalidURL)
	}
	if req.DemoURL != nil && !sanitize.ValidateHTTPSURL(*req.DemoURL) {
		return nil, fmt.Errorf("demo_url: %w", ErrInvalidURL)
	}

	// 3. Category Type validation
	if req.CategoryID != nil && *req.CategoryID != "" {
		if err := s.categorySvc.ValidateCategoryForType(ctx, *req.CategoryID, "project"); err != nil {
			return nil, err
		}
	} else {
		req.CategoryID = nil
	}

	// 4. Skills validation
	if len(req.SkillIDs) > 0 {
		if _, err := s.skillSvc.ValidateSkillIDs(ctx, req.SkillIDs); err != nil {
			return nil, err
		}
	}

	// 5. Date validation
	if req.StartedAt != nil && req.EndedAt != nil {
		if req.EndedAt.Before(*req.StartedAt) {
			return nil, ErrInvalidDateSpan
		}
	}

	// 6. Title, Slug & Markdown sanitization
	title := sanitize.NormalizeText(req.Title)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = GenerateSlug(title)
	} else {
		slug = GenerateSlug(slug)
	}

	// Conflict check on slug
	if existing, _ := s.repo.GetBySlug(ctx, slug); existing != nil {
		return nil, repository.ErrProjectConflict
	}

	cleanMD := sanitize.NormalizeText(req.DescriptionMD)
	html, err := sanitize.MarkdownToHTML(cleanMD)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown: %w", err)
	}

	project := &domain.Project{
		CategoryID:      req.CategoryID,
		Title:           title,
		Slug:            slug,
		Summary:         sanitize.NormalizeText(req.Summary),
		DescriptionMD:   cleanMD,
		DescriptionHTML: html,
		ThumbnailPath:   storage.SanitizeStoragePathPtr(req.ThumbnailPath),
		RepoURL:         req.RepoURL,
		DemoURL:         req.DemoURL,
		Status:          status,
		IsFeatured:      req.IsFeatured,
		SortOrder:       req.SortOrder,
		StartedAt:       req.StartedAt,
		EndedAt:         req.EndedAt,
	}

	if err := s.repo.Create(ctx, project, req.SkillIDs); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return s.repo.GetByID(ctx, project.ID)
}

func (s *projectService) UpdateProject(ctx context.Context, id string, req domain.ProjectUpdateRequest) (*domain.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status != "draft" && status != "published" && status != "archived" {
		return nil, ErrInvalidStatus
	}

	if req.RepoURL != nil && !sanitize.ValidateHTTPSURL(*req.RepoURL) {
		return nil, fmt.Errorf("repo_url: %w", ErrInvalidURL)
	}
	if req.DemoURL != nil && !sanitize.ValidateHTTPSURL(*req.DemoURL) {
		return nil, fmt.Errorf("demo_url: %w", ErrInvalidURL)
	}

	if req.CategoryID != nil && *req.CategoryID != "" {
		if err := s.categorySvc.ValidateCategoryForType(ctx, *req.CategoryID, "project"); err != nil {
			return nil, err
		}
	} else {
		req.CategoryID = nil
	}

	if len(req.SkillIDs) > 0 {
		if _, err := s.skillSvc.ValidateSkillIDs(ctx, req.SkillIDs); err != nil {
			return nil, err
		}
	}

	if req.StartedAt != nil && req.EndedAt != nil {
		if req.EndedAt.Before(*req.StartedAt) {
			return nil, ErrInvalidDateSpan
		}
	}

	title := sanitize.NormalizeText(req.Title)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = GenerateSlug(title)
	} else {
		slug = GenerateSlug(slug)
	}

	if slug != project.Slug {
		if existing, _ := s.repo.GetBySlug(ctx, slug); existing != nil && existing.ID != id {
			return nil, repository.ErrProjectConflict
		}
	}

	cleanMD := sanitize.NormalizeText(req.DescriptionMD)
	html, err := sanitize.MarkdownToHTML(cleanMD)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown: %w", err)
	}

	project.CategoryID = req.CategoryID
	project.Title = title
	project.Slug = slug
	project.Summary = sanitize.NormalizeText(req.Summary)
	project.DescriptionMD = cleanMD
	project.DescriptionHTML = html
	project.ThumbnailPath = storage.SanitizeStoragePathPtr(req.ThumbnailPath)
	project.RepoURL = req.RepoURL
	project.DemoURL = req.DemoURL
	project.Status = status
	project.IsFeatured = req.IsFeatured
	project.SortOrder = req.SortOrder
	project.StartedAt = req.StartedAt
	project.EndedAt = req.EndedAt

	if err := s.repo.Update(ctx, id, project, req.SkillIDs); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *projectService) UpdateProjectSkills(ctx context.Context, id string, skillIDs []string) error {
	if len(skillIDs) > 0 {
		if _, err := s.skillSvc.ValidateSkillIDs(ctx, skillIDs); err != nil {
			return err
		}
	}
	return s.repo.UpdateSkills(ctx, id, skillIDs)
}

func (s *projectService) DeleteProject(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
