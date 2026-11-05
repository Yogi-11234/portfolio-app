package domain

import "time"

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Type      string    `json:"type"` // "project" or "post"
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CategoryCreateRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=80"`
	Slug      string `json:"slug" validate:"omitempty,max=100"`
	Type      string `json:"type" validate:"required,oneof=project post"`
	SortOrder int    `json:"sort_order"`
}

type CategoryUpdateRequest struct {
	Name      string  `json:"name" validate:"required,min=2,max=80"`
	Slug      string  `json:"slug" validate:"omitempty,max=100"`
	Type      *string `json:"type" validate:"omitempty,oneof=project post"`
	SortOrder int     `json:"sort_order"`
}
