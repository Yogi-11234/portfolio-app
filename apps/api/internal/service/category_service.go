package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/sanitize"
)

var (
	ErrInvalidCategoryType = errors.New("tipe kategori harus 'project' atau 'post'")
	slugRegex              = regexp.MustCompile(`[^a-z0-9]+`)
)

func GenerateSlug(input string) string {
	s := strings.ToLower(sanitize.NormalizeText(input))
	s = slugRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

type CategoryService interface {
	ListCategories(ctx context.Context, categoryType string) ([]domain.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*domain.Category, error)
	CreateCategory(ctx context.Context, req domain.CategoryCreateRequest) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id string, req domain.CategoryUpdateRequest) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	ValidateCategoryForType(ctx context.Context, categoryID, expectedType string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) ListCategories(ctx context.Context, categoryType string) ([]domain.Category, error) {
	if categoryType != "" && categoryType != "project" && categoryType != "post" {
		return nil, ErrInvalidCategoryType
	}
	return s.repo.List(ctx, categoryType)
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *categoryService) CreateCategory(ctx context.Context, req domain.CategoryCreateRequest) (*domain.Category, error) {
	catType := strings.ToLower(strings.TrimSpace(req.Type))
	if catType != "project" && catType != "post" {
		return nil, ErrInvalidCategoryType
	}

	name := sanitize.NormalizeText(req.Name)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = GenerateSlug(name)
	} else {
		slug = GenerateSlug(slug)
	}

	// Check conflict
	existing, err := s.repo.GetBySlugAndType(ctx, slug, catType)
	if err == nil && existing != nil {
		return nil, repository.ErrCategoryConflict
	}

	category := &domain.Category{
		Name:      name,
		Slug:      slug,
		Type:      catType,
		SortOrder: req.SortOrder,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req domain.CategoryUpdateRequest) (*domain.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := sanitize.NormalizeText(req.Name)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = GenerateSlug(name)
	} else {
		slug = GenerateSlug(slug)
	}

	// If slug or type changed, check conflict within target type
	targetType := category.Type
	if req.Type != nil && *req.Type != "" {
		targetType = *req.Type
	}

	if slug != category.Slug || targetType != category.Type {
		existing, err := s.repo.GetBySlugAndType(ctx, slug, targetType)
		if err == nil && existing != nil && existing.ID != id {
			return nil, repository.ErrCategoryConflict
		}
	}

	category.Name = name
	category.Slug = slug
	category.Type = targetType
	category.SortOrder = req.SortOrder

	if err := s.repo.Update(ctx, id, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *categoryService) ValidateCategoryForType(ctx context.Context, categoryID, expectedType string) error {
	if categoryID == "" {
		return nil
	}
	cat, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("category_id: %w", repository.ErrCategoryNotFound)
	}
	if cat.Type != expectedType {
		return fmt.Errorf("kategori '%s' bukan untuk tipe '%s'", cat.Name, expectedType)
	}
	return nil
}
