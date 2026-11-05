package domain

import (
	"encoding/json"
	"time"

	"portfolio/api/internal/storage"
)

type Post struct {
	ID             string         `json:"id"`
	AuthorID       string         `json:"author_id"`
	CategoryID     *string        `json:"category_id"`
	Title          string         `json:"title"`
	Slug           string         `json:"slug"`
	Excerpt        string         `json:"excerpt"`
	ContentMD      string         `json:"content_md"`
	ContentHTML    string         `json:"content_html"`
	CoverPath      *string        `json:"cover_path"`
	Status         string         `json:"status"` // 'draft','published','archived'
	ReadingMinutes int16          `json:"reading_minutes"`
	PublishedAt    *time.Time     `json:"published_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Category       *CategoryBrief `json:"category,omitempty"`
	AuthorName     string         `json:"author_name,omitempty"`
}

func (p Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	return json.Marshal(&struct {
		Alias
		CoverPath *string `json:"cover_path"`
	}{
		Alias:     Alias(p),
		CoverPath: storage.FormatMediaURLPtr(p.CoverPath),
	})
}

type PostFilter struct {
	CategorySlug string
	Status       string // 'published' for public, or any for admin
	Page         int
	PerPage      int
}

type PostCreateRequest struct {
	Title          string  `json:"title" validate:"required,min=3,max=200"`
	Slug           string  `json:"slug" validate:"omitempty,max=220"`
	CategoryID     *string `json:"category_id" validate:"omitempty,uuid"`
	Excerpt        string  `json:"excerpt" validate:"omitempty,max=320"`
	ContentMD      string  `json:"content_md" validate:"required,min=5"`
	CoverPath      *string `json:"cover_path" validate:"omitempty,max=255"`
	Status         string  `json:"status" validate:"required,oneof=draft published archived"`
	ReadingMinutes *int16  `json:"reading_minutes" validate:"omitempty,min=1"`
}

type PostUpdateRequest struct {
	Title          string  `json:"title" validate:"required,min=3,max=200"`
	Slug           string  `json:"slug" validate:"omitempty,max=220"`
	CategoryID     *string `json:"category_id" validate:"omitempty,uuid"`
	Excerpt        string  `json:"excerpt" validate:"omitempty,max=320"`
	ContentMD      string  `json:"content_md" validate:"required,min=5"`
	CoverPath      *string `json:"cover_path" validate:"omitempty,max=255"`
	Status         string  `json:"status" validate:"required,oneof=draft published archived"`
	ReadingMinutes *int16  `json:"reading_minutes" validate:"omitempty,min=1"`
}
