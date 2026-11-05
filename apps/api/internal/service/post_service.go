package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/sanitize"
	"portfolio/api/internal/storage"
)

type PostService interface {
	ListPosts(ctx context.Context, f domain.PostFilter) ([]domain.Post, int, error)
	GetPostByID(ctx context.Context, id string) (*domain.Post, error)
	GetPostBySlug(ctx context.Context, slug string) (*domain.Post, error)
	CreatePost(ctx context.Context, authorID string, req domain.PostCreateRequest) (*domain.Post, error)
	UpdatePost(ctx context.Context, id string, req domain.PostUpdateRequest) (*domain.Post, error)
	DeletePost(ctx context.Context, id string) error
}

type postService struct {
	repo        repository.PostRepository
	categorySvc CategoryService
}

func NewPostService(repo repository.PostRepository, categorySvc CategoryService) PostService {
	return &postService{
		repo:        repo,
		categorySvc: categorySvc,
	}
}

func (s *postService) ListPosts(ctx context.Context, f domain.PostFilter) ([]domain.Post, int, error) {
	return s.repo.List(ctx, f)
}

func (s *postService) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *postService) GetPostBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	cleanSlug := strings.ToLower(strings.TrimSpace(slug))
	return s.repo.GetBySlug(ctx, cleanSlug)
}

func calculateReadingMinutes(content string) int16 {
	words := len(strings.Fields(content))
	minutes := int16(words / 200)
	if minutes < 1 {
		return 1
	}
	return minutes
}

func (s *postService) CreatePost(ctx context.Context, authorID string, req domain.PostCreateRequest) (*domain.Post, error) {
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status != "draft" && status != "published" && status != "archived" {
		return nil, ErrInvalidStatus
	}

	if req.CategoryID != nil && *req.CategoryID != "" {
		if err := s.categorySvc.ValidateCategoryForType(ctx, *req.CategoryID, "post"); err != nil {
			return nil, err
		}
	} else {
		req.CategoryID = nil
	}

	title := sanitize.NormalizeText(req.Title)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = GenerateSlug(title)
	} else {
		slug = GenerateSlug(slug)
	}

	if existing, _ := s.repo.GetBySlug(ctx, slug); existing != nil {
		return nil, repository.ErrPostConflict
	}

	cleanMD := sanitize.NormalizeText(req.ContentMD)
	html, err := sanitize.MarkdownToHTML(cleanMD)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown: %w", err)
	}

	excerpt := sanitize.NormalizeText(req.Excerpt)
	if excerpt == "" {
		if len(cleanMD) > 200 {
			excerpt = cleanMD[:200] + "..."
		} else {
			excerpt = cleanMD
		}
	}

	readingMinutes := calculateReadingMinutes(cleanMD)
	if req.ReadingMinutes != nil && *req.ReadingMinutes > 0 {
		readingMinutes = *req.ReadingMinutes
	}

	var publishedAt *time.Time
	if status == "published" {
		now := time.Now()
		publishedAt = &now
	}

	post := &domain.Post{
		AuthorID:       authorID,
		CategoryID:     req.CategoryID,
		Title:          title,
		Slug:           slug,
		Excerpt:        excerpt,
		ContentMD:      cleanMD,
		ContentHTML:    html,
		CoverPath:      storage.SanitizeStoragePathPtr(req.CoverPath),
		Status:         status,
		ReadingMinutes: readingMinutes,
		PublishedAt:    publishedAt,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return s.repo.GetByID(ctx, post.ID)
}

func (s *postService) UpdatePost(ctx context.Context, id string, req domain.PostUpdateRequest) (*domain.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status != "draft" && status != "published" && status != "archived" {
		return nil, ErrInvalidStatus
	}

	if req.CategoryID != nil && *req.CategoryID != "" {
		if err := s.categorySvc.ValidateCategoryForType(ctx, *req.CategoryID, "post"); err != nil {
			return nil, err
		}
	} else {
		req.CategoryID = nil
	}

	title := sanitize.NormalizeText(req.Title)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = GenerateSlug(title)
	} else {
		slug = GenerateSlug(slug)
	}

	if slug != post.Slug {
		if existing, _ := s.repo.GetBySlug(ctx, slug); existing != nil && existing.ID != id {
			return nil, repository.ErrPostConflict
		}
	}

	cleanMD := sanitize.NormalizeText(req.ContentMD)
	html, err := sanitize.MarkdownToHTML(cleanMD)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown: %w", err)
	}

	excerpt := sanitize.NormalizeText(req.Excerpt)
	if excerpt == "" {
		if len(cleanMD) > 200 {
			excerpt = cleanMD[:200] + "..."
		} else {
			excerpt = cleanMD
		}
	}

	readingMinutes := calculateReadingMinutes(cleanMD)
	if req.ReadingMinutes != nil && *req.ReadingMinutes > 0 {
		readingMinutes = *req.ReadingMinutes
	}

	publishedAt := post.PublishedAt
	if status == "published" && publishedAt == nil {
		now := time.Now()
		publishedAt = &now
	}

	post.CategoryID = req.CategoryID
	post.Title = title
	post.Slug = slug
	post.Excerpt = excerpt
	post.ContentMD = cleanMD
	post.ContentHTML = html
	post.CoverPath = storage.SanitizeStoragePathPtr(req.CoverPath)
	post.Status = status
	post.ReadingMinutes = readingMinutes
	post.PublishedAt = publishedAt

	if err := s.repo.Update(ctx, id, post); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *postService) DeletePost(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
