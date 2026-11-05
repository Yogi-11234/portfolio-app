package service

import (
	"context"
	"testing"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
)

type mockCategoryRepo struct {
	categories []domain.Category
}

func (m *mockCategoryRepo) List(ctx context.Context, categoryType string) ([]domain.Category, error) {
	var result []domain.Category
	for _, c := range m.categories {
		if categoryType == "" || c.Type == categoryType {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	for _, c := range m.categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, repository.ErrCategoryNotFound
}

func (m *mockCategoryRepo) GetBySlugAndType(ctx context.Context, slug, categoryType string) (*domain.Category, error) {
	for _, c := range m.categories {
		if c.Slug == slug && c.Type == categoryType {
			return &c, nil
		}
	}
	return nil, repository.ErrCategoryNotFound
}

func (m *mockCategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	c.ID = "cat-uuid-" + c.Slug
	m.categories = append(m.categories, *c)
	return nil
}

func (m *mockCategoryRepo) Update(ctx context.Context, id string, c *domain.Category) error {
	for i, existing := range m.categories {
		if existing.ID == id {
			m.categories[i] = *c
			return nil
		}
	}
	return repository.ErrCategoryNotFound
}

func (m *mockCategoryRepo) Delete(ctx context.Context, id string) error {
	for i, existing := range m.categories {
		if existing.ID == id {
			m.categories = append(m.categories[:i], m.categories[i+1:]...)
			return nil
		}
	}
	return repository.ErrCategoryNotFound
}

func TestCategoryService(t *testing.T) {
	mockRepo := &mockCategoryRepo{
		categories: []domain.Category{
			{ID: "cat-1", Name: "Web Apps", Slug: "web-apps", Type: "project"},
		},
	}
	svc := NewCategoryService(mockRepo)
	ctx := context.Background()

	// 1. Invalid category type
	_, err := svc.CreateCategory(ctx, domain.CategoryCreateRequest{
		Name: "Invalid",
		Type: "unknown",
	})
	if err != ErrInvalidCategoryType {
		t.Errorf("Expected ErrInvalidCategoryType, got: %v", err)
	}

	// 2. Conflict on same slug and type
	_, err = svc.CreateCategory(ctx, domain.CategoryCreateRequest{
		Name: "Web Apps",
		Type: "project",
	})
	if err != repository.ErrCategoryConflict {
		t.Errorf("Expected ErrCategoryConflict, got: %v", err)
	}

	// 3. Same slug on DIFFERENT type is allowed
	postCat, err := svc.CreateCategory(ctx, domain.CategoryCreateRequest{
		Name: "Web Apps",
		Type: "post",
	})
	if err != nil {
		t.Fatalf("Expected create on different type to succeed, got: %v", err)
	}
	if postCat.Type != "post" || postCat.Slug != "web-apps" {
		t.Errorf("Unexpected created category: %+v", postCat)
	}

	// 4. Validate category type rule (post cat cannot be used for project)
	err = svc.ValidateCategoryForType(ctx, postCat.ID, "project")
	if err == nil {
		t.Error("Expected error validating post category for project")
	}

	err = svc.ValidateCategoryForType(ctx, postCat.ID, "post")
	if err != nil {
		t.Errorf("Expected valid type match, got: %v", err)
	}
}
