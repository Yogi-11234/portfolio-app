package domain

import (
	"encoding/json"
	"time"

	"portfolio/api/internal/storage"
)

type SocialLink struct {
	Label string `json:"label" validate:"required,max=50"`
	URL   string `json:"url" validate:"required,url"`
}

type Profile struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	FullName    string          `json:"full_name"`
	Headline    string          `json:"headline"`
	BioMD       string          `json:"bio_md"`
	BioHTML     string          `json:"bio_html"`
	Location    *string         `json:"location"`
	PublicEmail *string         `json:"public_email"`
	AvatarPath  *string         `json:"avatar_path"`
	ResumePath  *string         `json:"resume_path"`
	GithubURL   *string         `json:"github_url"`
	LinkedinURL *string         `json:"linkedin_url"`
	SocialLinks json.RawMessage `json:"social_links"`
	OpenToWork  bool            `json:"open_to_work"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (p Profile) MarshalJSON() ([]byte, error) {
	type Alias Profile
	return json.Marshal(&struct {
		Alias
		AvatarPath *string `json:"avatar_path"`
		ResumePath *string `json:"resume_path"`
	}{
		Alias:      Alias(p),
		AvatarPath: storage.FormatMediaURLPtr(p.AvatarPath),
		ResumePath: storage.FormatMediaURLPtr(p.ResumePath),
	})
}

type ProfileUpdateRequest struct {
	FullName    string       `json:"full_name" validate:"required,min=2,max=120"`
	Headline    string       `json:"headline" validate:"required,min=2,max=160"`
	BioMD       string       `json:"bio_md" validate:"max=50000"`
	Location    *string      `json:"location" validate:"omitempty,max=120"`
	PublicEmail *string      `json:"public_email" validate:"omitempty,email,max=160"`
	AvatarPath  *string      `json:"avatar_path" validate:"omitempty,max=255"`
	ResumePath  *string      `json:"resume_path" validate:"omitempty,max=255"`
	GithubURL   *string      `json:"github_url" validate:"omitempty,max=255"`
	LinkedinURL *string      `json:"linkedin_url" validate:"omitempty,max=255"`
	SocialLinks []SocialLink `json:"social_links" validate:"dive"`
	OpenToWork  bool         `json:"open_to_work"`
}

// UnmarshalJSON provides flexible decoding of SocialLinks (accepting both array or map object)
// to prevent "invalid JSON" errors when client sends social links as key-value pairs.
func (r *ProfileUpdateRequest) UnmarshalJSON(data []byte) error {
	type Alias ProfileUpdateRequest
	var raw struct {
		Alias
		RawSocialLinks json.RawMessage `json:"social_links"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = ProfileUpdateRequest(raw.Alias)

	if len(raw.RawSocialLinks) > 0 {
		// 1. Try decoding as slice of SocialLink
		var list []SocialLink
		if err := json.Unmarshal(raw.RawSocialLinks, &list); err == nil {
			r.SocialLinks = list
		} else {
			// 2. Try decoding as map[string]string (e.g. {"github": "https://..."})
			var m map[string]string
			if errMap := json.Unmarshal(raw.RawSocialLinks, &m); errMap == nil {
				r.SocialLinks = make([]SocialLink, 0, len(m))
				for k, v := range m {
					if v != "" {
						r.SocialLinks = append(r.SocialLinks, SocialLink{
							Label: k,
							URL:   v,
						})
					}
				}
			}
		}
	}
	return nil
}
