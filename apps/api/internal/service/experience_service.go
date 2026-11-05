package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/sanitize"
)

var (
	ErrInvalidExperienceType = errors.New("tipe pengalaman harus salah satu dari: work, education, organization, volunteer")
)

var allowedExperienceTypes = map[string]bool{
	"work":         true,
	"education":    true,
	"organization": true,
	"volunteer":    true,
}

type ExperienceService interface {
	ListExperiences(ctx context.Context, onlyVisible bool, expType string) ([]domain.Experience, error)
	GetExperienceByID(ctx context.Context, id string) (*domain.Experience, error)
	CreateExperience(ctx context.Context, req domain.ExperienceCreateRequest) (*domain.Experience, error)
	UpdateExperience(ctx context.Context, id string, req domain.ExperienceUpdateRequest) (*domain.Experience, error)
	DeleteExperience(ctx context.Context, id string) error
}

type experienceService struct {
	repo repository.ExperienceRepository
}

func NewExperienceService(repo repository.ExperienceRepository) ExperienceService {
	return &experienceService{repo: repo}
}

func (s *experienceService) ListExperiences(ctx context.Context, onlyVisible bool, expType string) ([]domain.Experience, error) {
	cleanType := strings.ToLower(strings.TrimSpace(expType))
	if cleanType != "" && !allowedExperienceTypes[cleanType] {
		return nil, ErrInvalidExperienceType
	}
	return s.repo.List(ctx, onlyVisible, cleanType)
}

func (s *experienceService) GetExperienceByID(ctx context.Context, id string) (*domain.Experience, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *experienceService) CreateExperience(ctx context.Context, req domain.ExperienceCreateRequest) (*domain.Experience, error) {
	expType := strings.ToLower(strings.TrimSpace(req.Type))
	if !allowedExperienceTypes[expType] {
		return nil, ErrInvalidExperienceType
	}

	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDateSpan
	}

	cleanMD := sanitize.NormalizeText(req.DescriptionMD)
	html, err := sanitize.MarkdownToHTML(cleanMD)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown: %w", err)
	}

	isVisible := true
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}

	exp := &domain.Experience{
		Type:            expType,
		Organization:    sanitize.NormalizeText(req.Organization),
		Role:            sanitize.NormalizeText(req.Role),
		Location:        req.Location,
		DescriptionMD:   cleanMD,
		DescriptionHTML: html,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		IsCurrent:       req.IsCurrent,
		IsVisible:       isVisible,
		SortOrder:       req.SortOrder,
	}

	if err := s.repo.Create(ctx, exp); err != nil {
		return nil, fmt.Errorf("failed to create experience: %w", err)
	}

	return exp, nil
}

func (s *experienceService) UpdateExperience(ctx context.Context, id string, req domain.ExperienceUpdateRequest) (*domain.Experience, error) {
	exp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	expType := strings.ToLower(strings.TrimSpace(req.Type))
	if !allowedExperienceTypes[expType] {
		return nil, ErrInvalidExperienceType
	}

	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDateSpan
	}

	cleanMD := sanitize.NormalizeText(req.DescriptionMD)
	html, err := sanitize.MarkdownToHTML(cleanMD)
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown: %w", err)
	}

	isVisible := exp.IsVisible
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}

	exp.Type = expType
	exp.Organization = sanitize.NormalizeText(req.Organization)
	exp.Role = sanitize.NormalizeText(req.Role)
	exp.Location = req.Location
	exp.DescriptionMD = cleanMD
	exp.DescriptionHTML = html
	exp.StartDate = req.StartDate
	exp.EndDate = req.EndDate
	exp.IsCurrent = req.IsCurrent
	exp.IsVisible = isVisible
	exp.SortOrder = req.SortOrder

	if err := s.repo.Update(ctx, id, exp); err != nil {
		return nil, err
	}

	return exp, nil
}

func (s *experienceService) DeleteExperience(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
