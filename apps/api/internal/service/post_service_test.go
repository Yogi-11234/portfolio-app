package service

import (
	"context"
	"strings"
	"testing"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
)

type mockPostRepo struct {
	posts []domain.Post
}

func (m *mockPostRepo) List(ctx context.Context, f domain.PostFilter) ([]domain.Post, int, error) {
	return m.posts, len(m.posts), nil
}

func (m *mockPostRepo) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	for _, p := range m.posts {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, repository.ErrPostNotFound
}

func (m *mockPostRepo) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	for _, p := range m.posts {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, repository.ErrPostNotFound
}

func (m *mockPostRepo) Create(ctx context.Context, p *domain.Post) error {
	p.ID = "post-uuid-" + p.Slug
	m.posts = append(m.posts, *p)
	return nil
}

func (m *mockPostRepo) Update(ctx context.Context, id string, p *domain.Post) error {
	for i, existing := range m.posts {
		if existing.ID == id {
			m.posts[i] = *p
			return nil
		}
	}
	return repository.ErrPostNotFound
}

func (m *mockPostRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func TestPostService(t *testing.T) {
	mockRepo := &mockPostRepo{}
	mockCatRepo := &mockCategoryRepo{
		categories: []domain.Category{
			{ID: "cat-post", Name: "Engineering", Slug: "engineering", Type: "post"},
			{ID: "cat-proj", Name: "Projects", Slug: "projects", Type: "project"},
		},
	}
	catSvc := NewCategoryService(mockCatRepo)
	svc := NewPostService(mockRepo, catSvc)
	ctx := context.Background()

	// 1. Invalid status
	_, err := svc.CreatePost(ctx, "author-1", domain.PostCreateRequest{
		Title:     "My First Post",
		ContentMD: "Some long content...",
		Status:    "invalid_status",
	})
	if err != ErrInvalidStatus {
		t.Errorf("Expected ErrInvalidStatus, got: %v", err)
	}

	// 2. Category mismatch (using project category for blog post)
	projCatID := "cat-proj"
	_, err = svc.CreatePost(ctx, "author-1", domain.PostCreateRequest{
		Title:      "My First Post",
		ContentMD:  "Some long content...",
		Status:     "published",
		CategoryID: &projCatID,
	})
	if err == nil || !strings.Contains(err.Error(), "bukan untuk tipe") {
		t.Errorf("Expected category mismatch error, got: %v", err)
	}

	// 3. Valid Post
	postCatID := "cat-post"
	post, err := svc.CreatePost(ctx, "author-1", domain.PostCreateRequest{
		Title:      "Modern Go Architecture",
		ContentMD:  "# Modern Go\n\nDesigning **clean** systems in Go with <script>evil()</script>",
		Status:     "published",
		CategoryID: &postCatID,
	})
	if err != nil {
		t.Fatalf("Expected create post to succeed, got: %v", err)
	}
	if post.Slug != "modern-go-architecture" {
		t.Errorf("Expected slug 'modern-go-architecture', got '%s'", post.Slug)
	}
	if post.PublishedAt == nil {
		t.Error("Expected PublishedAt to be set for published post")
	}
	if strings.Contains(post.ContentHTML, "<script") {
		t.Error("Expected ContentHTML to be sanitized")
	}
}
