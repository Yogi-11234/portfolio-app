package domain

import (
	"encoding/json"
	"time"

	"portfolio/api/internal/storage"
)

type Project struct {
	ID              string        `json:"id"`
	CategoryID      *string       `json:"category_id"`
	Title           string        `json:"title"`
	Slug            string        `json:"slug"`
	Summary         string        `json:"summary"`
	DescriptionMD   string        `json:"description_md"`
	DescriptionHTML string        `json:"description_html"`
	ThumbnailPath   *string       `json:"thumbnail_path"`
	DemoURL         *string       `json:"demo_url"`
	RepoURL         *string       `json:"repo_url"`
	Status          string        `json:"status"` // draft, published, archived
	IsFeatured      bool          `json:"is_featured"`
	SortOrder       int           `json:"sort_order"`
	StartedAt       *time.Time    `json:"started_at"`
	EndedAt         *time.Time    `json:"ended_at"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	Category        *CategoryBrief `json:"category,omitempty"`
	Skills          []SkillBrief  `json:"skills"`
}

func (p Project) MarshalJSON() ([]byte, error) {
	type Alias Project
	return json.Marshal(&struct {
		Alias
		ThumbnailPath *string `json:"thumbnail_path"`
	}{
		Alias:         Alias(p),
		ThumbnailPath: storage.FormatMediaURLPtr(p.ThumbnailPath),
	})
}

type ProjectFilter struct {
	CategorySlug string
	SkillSlug    string
	FeaturedOnly *bool
	Status       string // If empty, defaults to 'published' for public
	Page         int
	PerPage      int
}

type ProjectCreateRequest struct {
	Title         string     `json:"title" validate:"required,min=3,max=160"`
	Slug          string     `json:"slug" validate:"omitempty,max=180"`
	Summary       string     `json:"summary" validate:"required,max=300"`
	DescriptionMD string     `json:"description_md" validate:"max=50000"`
	CategoryID    *string    `json:"category_id" validate:"omitempty,uuid"`
	ThumbnailPath *string    `json:"thumbnail_path" validate:"omitempty,max=255"`
	RepoURL       *string    `json:"repo_url" validate:"omitempty,max=255"`
	DemoURL       *string    `json:"demo_url" validate:"omitempty,max=255"`
	Status        string     `json:"status" validate:"required,oneof=draft published archived"`
	IsFeatured    bool       `json:"is_featured"`
	SortOrder     int        `json:"sort_order"`
	StartedAt     *time.Time `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
	SkillIDs      []string   `json:"skill_ids" validate:"dive,uuid"`
}

func (r *ProjectCreateRequest) UnmarshalJSON(data []byte) error {
	type Alias ProjectCreateRequest
	var raw struct {
		Alias
		Featured   *bool      `json:"featured"`
		FinishedAt *time.Time `json:"finished_at"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = ProjectCreateRequest(raw.Alias)
	if !r.IsFeatured && raw.Featured != nil {
		r.IsFeatured = *raw.Featured
	}
	if r.EndedAt == nil && raw.FinishedAt != nil {
		r.EndedAt = raw.FinishedAt
	}
	return nil
}

type ProjectUpdateRequest struct {
	Title         string     `json:"title" validate:"required,min=3,max=160"`
	Slug          string     `json:"slug" validate:"omitempty,max=180"`
	Summary       string     `json:"summary" validate:"required,max=300"`
	DescriptionMD string     `json:"description_md" validate:"max=50000"`
	CategoryID    *string    `json:"category_id" validate:"omitempty,uuid"`
	ThumbnailPath *string    `json:"thumbnail_path" validate:"omitempty,max=255"`
	RepoURL       *string    `json:"repo_url" validate:"omitempty,max=255"`
	DemoURL       *string    `json:"demo_url" validate:"omitempty,max=255"`
	Status        string     `json:"status" validate:"required,oneof=draft published archived"`
	IsFeatured    bool       `json:"is_featured"`
	SortOrder     int        `json:"sort_order"`
	StartedAt     *time.Time `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
	SkillIDs      []string   `json:"skill_ids" validate:"dive,uuid"`
}

func (r *ProjectUpdateRequest) UnmarshalJSON(data []byte) error {
	type Alias ProjectUpdateRequest
	var raw struct {
		Alias
		Featured   *bool      `json:"featured"`
		FinishedAt *time.Time `json:"finished_at"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = ProjectUpdateRequest(raw.Alias)
	if !r.IsFeatured && raw.Featured != nil {
		r.IsFeatured = *raw.Featured
	}
	if r.EndedAt == nil && raw.FinishedAt != nil {
		r.EndedAt = raw.FinishedAt
	}
	return nil
}

type ProjectSkillsUpdateRequest struct {
	SkillIDs []string `json:"skill_ids" validate:"dive,uuid"`
}
