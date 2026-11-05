package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/sanitize"
	"portfolio/api/internal/storage"
)

var (
	ErrInvalidURL = errors.New("URL eksternal harus diawali dengan https:// dan tidak boleh mengandung skrip berbahaya")
)

type ProfileService interface {
	GetProfile(ctx context.Context) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, userID string, req domain.ProfileUpdateRequest) (*domain.Profile, error)
}

type profileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) GetProfile(ctx context.Context) (*domain.Profile, error) {
	return s.repo.Get(ctx)
}

func (s *profileService) UpdateProfile(ctx context.Context, userID string, req domain.ProfileUpdateRequest) (*domain.Profile, error) {
	// 1. Validate URLs
	if req.GithubURL != nil && !sanitize.ValidateHTTPSURL(*req.GithubURL) {
		return nil, fmt.Errorf("github_url: %w", ErrInvalidURL)
	}
	if req.LinkedinURL != nil && !sanitize.ValidateHTTPSURL(*req.LinkedinURL) {
		return nil, fmt.Errorf("linkedin_url: %w", ErrInvalidURL)
	}
	for _, sl := range req.SocialLinks {
		if !sanitize.ValidateHTTPSURL(sl.URL) {
			return nil, fmt.Errorf("social_links (%s): %w", sl.Label, ErrInvalidURL)
		}
	}

	// 2. Render & Sanitize Bio Markdown
	cleanBioMD := sanitize.NormalizeText(req.BioMD)
	bioHTML, err := sanitize.MarkdownToHTML(cleanBioMD)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize markdown: %w", err)
	}

	// 3. Serialize Social Links
	socialLinksBytes, err := json.Marshal(req.SocialLinks)
	if err != nil {
		socialLinksBytes = []byte("[]")
	}

	profile := &domain.Profile{
		FullName:    sanitize.NormalizeText(req.FullName),
		Headline:    sanitize.NormalizeText(req.Headline),
		BioMD:       cleanBioMD,
		BioHTML:     bioHTML,
		Location:    req.Location,
		PublicEmail: req.PublicEmail,
		AvatarPath:  storage.SanitizeStoragePathPtr(req.AvatarPath),
		ResumePath:  storage.SanitizeStoragePathPtr(req.ResumePath),
		GithubURL:   req.GithubURL,
		LinkedinURL: req.LinkedinURL,
		SocialLinks: json.RawMessage(socialLinksBytes),
		OpenToWork:  req.OpenToWork,
	}

	if err := s.repo.Update(ctx, userID, profile); err != nil {
		return nil, err
	}

	return s.repo.GetByUserID(ctx, userID)
}
