package domain

import (
	"encoding/json"
	"strings"

	"portfolio/api/internal/storage"
)

type Skill struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	GroupName  string  `json:"group_name"`
	IconKey    *string `json:"icon_key"`
	BrandColor *string `json:"brand_color"`
	Level      *int16  `json:"level"`
	IsVisible  bool    `json:"is_visible"`
	SortOrder  int     `json:"sort_order"`
}

func (s Skill) MarshalJSON() ([]byte, error) {
	type Alias Skill
	orderIndex := s.SortOrder
	iconName := ""
	if s.IconKey != nil {
		iconName = *s.IconKey
	}
	formattedIcon := s.IconKey
	if s.IconKey != nil && (strings.HasPrefix(*s.IconKey, "skills/") || strings.HasPrefix(*s.IconKey, "http") || strings.HasPrefix(*s.IconKey, "uploads/")) {
		f := storage.FormatMediaURL(*s.IconKey)
		formattedIcon = &f
		iconName = f
	}
	return json.Marshal(&struct {
		Alias
		OrderIndex int     `json:"order_index"`
		IconName   string  `json:"icon_name"`
		IconKey    *string `json:"icon_key"`
	}{
		Alias:      Alias(s),
		OrderIndex: orderIndex,
		IconName:   iconName,
		IconKey:    formattedIcon,
	})
}

type SkillBrief struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	IconKey    *string `json:"icon_key"`
	BrandColor *string `json:"brand_color"`
}

type SkillCreateRequest struct {
	Name       string  `json:"name" validate:"required,min=1,max=60"`
	Slug       string  `json:"slug" validate:"omitempty,max=80"`
	GroupName  string  `json:"group_name" validate:"required,oneof=language framework database tool cloud other"`
	IconKey    *string `json:"icon_key" validate:"omitempty,max=60"`
	BrandColor *string `json:"brand_color" validate:"omitempty"`
	Level      *int16  `json:"level" validate:"omitempty,min=1,max=5"`
	IsVisible  *bool   `json:"is_visible"`
	SortOrder  int     `json:"sort_order"`
}

func (r *SkillCreateRequest) UnmarshalJSON(data []byte) error {
	type Alias SkillCreateRequest
	var raw struct {
		Alias
		OrderIndex *int    `json:"order_index"`
		IconName   *string `json:"icon_name"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = SkillCreateRequest(raw.Alias)
	if r.SortOrder == 0 && raw.OrderIndex != nil {
		r.SortOrder = *raw.OrderIndex
	}
	if (r.IconKey == nil || *r.IconKey == "") && raw.IconName != nil && *raw.IconName != "" {
		r.IconKey = raw.IconName
	}
	return nil
}

type SkillUpdateRequest struct {
	Name       string  `json:"name" validate:"required,min=1,max=60"`
	Slug       string  `json:"slug" validate:"omitempty,max=80"`
	GroupName  string  `json:"group_name" validate:"required,oneof=language framework database tool cloud other"`
	IconKey    *string `json:"icon_key" validate:"omitempty,max=60"`
	BrandColor *string `json:"brand_color" validate:"omitempty"`
	Level      *int16  `json:"level" validate:"omitempty,min=1,max=5"`
	IsVisible  *bool   `json:"is_visible"`
	SortOrder  int     `json:"sort_order"`
}

func (r *SkillUpdateRequest) UnmarshalJSON(data []byte) error {
	type Alias SkillUpdateRequest
	var raw struct {
		Alias
		OrderIndex *int    `json:"order_index"`
		IconName   *string `json:"icon_name"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = SkillUpdateRequest(raw.Alias)
	if r.SortOrder == 0 && raw.OrderIndex != nil {
		r.SortOrder = *raw.OrderIndex
	}
	if (r.IconKey == nil || *r.IconKey == "") && raw.IconName != nil && *raw.IconName != "" {
		r.IconKey = raw.IconName
	}
	return nil
}
