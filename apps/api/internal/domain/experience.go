package domain

import (
	"encoding/json"
	"strings"
	"time"
)

type Experience struct {
	ID              string     `json:"id"`
	Type            string     `json:"type"` // work, education, organization, volunteer
	Organization    string     `json:"organization"`
	Role            string     `json:"role"`
	Location        *string    `json:"location"`
	DescriptionMD   string     `json:"description_md"`
	DescriptionHTML string     `json:"description_html"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	IsCurrent       bool       `json:"is_current"`
	IsVisible       bool       `json:"is_visible"`
	SortOrder       int        `json:"sort_order"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (e Experience) MarshalJSON() ([]byte, error) {
	type Alias Experience
	var endDateStr *string
	if e.EndDate != nil {
		s := e.EndDate.Format(time.RFC3339)
		endDateStr = &s
	}
	return json.Marshal(&struct {
		Alias
		Title      string  `json:"title"`
		StartedAt  string  `json:"started_at"`
		EndedAt    *string `json:"ended_at"`
		Current    bool    `json:"current"`
		OrderIndex int     `json:"order_index"`
	}{
		Alias:      Alias(e),
		Title:      e.Role,
		StartedAt:  e.StartDate.Format(time.RFC3339),
		EndedAt:    endDateStr,
		Current:    e.IsCurrent,
		OrderIndex: e.SortOrder,
	})
}

func parseFlexibleDate(val interface{}) (*time.Time, error) {
	if val == nil {
		return nil, nil
	}
	s, ok := val.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return nil, nil
	}
	s = strings.TrimSpace(s)
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return &t, nil
		}
	}
	return nil, nil
}

type ExperienceCreateRequest struct {
	Type          string     `json:"type" validate:"required,oneof=work education organization volunteer"`
	Organization  string     `json:"organization" validate:"required,min=2,max=160"`
	Role          string     `json:"role" validate:"required,min=2,max=160"`
	Location      *string    `json:"location" validate:"omitempty,max=120"`
	DescriptionMD string     `json:"description_md" validate:"max=50000"`
	StartDate     time.Time  `json:"start_date" validate:"required"`
	EndDate       *time.Time `json:"end_date"`
	IsCurrent     bool       `json:"is_current"`
	IsVisible     *bool      `json:"is_visible"`
	SortOrder     int        `json:"sort_order"`
}

func (r *ExperienceCreateRequest) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type          string      `json:"type"`
		Organization  string      `json:"organization"`
		Role          string      `json:"role"`
		Title         string      `json:"title"`
		Location      *string     `json:"location"`
		DescriptionMD string      `json:"description_md"`
		StartDate     interface{} `json:"start_date"`
		StartedAt     interface{} `json:"started_at"`
		EndDate       interface{} `json:"end_date"`
		EndedAt       interface{} `json:"ended_at"`
		IsCurrent     *bool       `json:"is_current"`
		Current       *bool       `json:"current"`
		IsVisible     *bool       `json:"is_visible"`
		SortOrder     *int        `json:"sort_order"`
		OrderIndex    *int        `json:"order_index"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Type = raw.Type
	r.Organization = raw.Organization
	if raw.Role != "" {
		r.Role = raw.Role
	} else {
		r.Role = raw.Title
	}
	r.Location = raw.Location
	r.DescriptionMD = raw.DescriptionMD

	var err error
	var sDate *time.Time
	if raw.StartDate != nil {
		sDate, err = parseFlexibleDate(raw.StartDate)
	} else if raw.StartedAt != nil {
		sDate, err = parseFlexibleDate(raw.StartedAt)
	}
	if err == nil && sDate != nil {
		r.StartDate = *sDate
	}

	var eDate *time.Time
	if raw.EndDate != nil {
		eDate, _ = parseFlexibleDate(raw.EndDate)
	} else if raw.EndedAt != nil {
		eDate, _ = parseFlexibleDate(raw.EndedAt)
	}
	r.EndDate = eDate

	if raw.IsCurrent != nil {
		r.IsCurrent = *raw.IsCurrent
	} else if raw.Current != nil {
		r.IsCurrent = *raw.Current
	}

	r.IsVisible = raw.IsVisible

	if raw.SortOrder != nil {
		r.SortOrder = *raw.SortOrder
	} else if raw.OrderIndex != nil {
		r.SortOrder = *raw.OrderIndex
	}

	return nil
}

type ExperienceUpdateRequest struct {
	Type          string     `json:"type" validate:"required,oneof=work education organization volunteer"`
	Organization  string     `json:"organization" validate:"required,min=2,max=160"`
	Role          string     `json:"role" validate:"required,min=2,max=160"`
	Location      *string    `json:"location" validate:"omitempty,max=120"`
	DescriptionMD string     `json:"description_md" validate:"max=50000"`
	StartDate     time.Time  `json:"start_date" validate:"required"`
	EndDate       *time.Time `json:"end_date"`
	IsCurrent     bool       `json:"is_current"`
	IsVisible     *bool      `json:"is_visible"`
	SortOrder     int        `json:"sort_order"`
}

func (r *ExperienceUpdateRequest) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type          string      `json:"type"`
		Organization  string      `json:"organization"`
		Role          string      `json:"role"`
		Title         string      `json:"title"`
		Location      *string     `json:"location"`
		DescriptionMD string      `json:"description_md"`
		StartDate     interface{} `json:"start_date"`
		StartedAt     interface{} `json:"started_at"`
		EndDate       interface{} `json:"end_date"`
		EndedAt       interface{} `json:"ended_at"`
		IsCurrent     *bool       `json:"is_current"`
		Current       *bool       `json:"current"`
		IsVisible     *bool       `json:"is_visible"`
		SortOrder     *int        `json:"sort_order"`
		OrderIndex    *int        `json:"order_index"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Type = raw.Type
	r.Organization = raw.Organization
	if raw.Role != "" {
		r.Role = raw.Role
	} else {
		r.Role = raw.Title
	}
	r.Location = raw.Location
	r.DescriptionMD = raw.DescriptionMD

	var err error
	var sDate *time.Time
	if raw.StartDate != nil {
		sDate, err = parseFlexibleDate(raw.StartDate)
	} else if raw.StartedAt != nil {
		sDate, err = parseFlexibleDate(raw.StartedAt)
	}
	if err == nil && sDate != nil {
		r.StartDate = *sDate
	}

	var eDate *time.Time
	if raw.EndDate != nil {
		eDate, _ = parseFlexibleDate(raw.EndDate)
	} else if raw.EndedAt != nil {
		eDate, _ = parseFlexibleDate(raw.EndedAt)
	}
	r.EndDate = eDate

	if raw.IsCurrent != nil {
		r.IsCurrent = *raw.IsCurrent
	} else if raw.Current != nil {
		r.IsCurrent = *raw.Current
	}

	r.IsVisible = raw.IsVisible

	if raw.SortOrder != nil {
		r.SortOrder = *raw.SortOrder
	} else if raw.OrderIndex != nil {
		r.SortOrder = *raw.OrderIndex
	}

	return nil
}
