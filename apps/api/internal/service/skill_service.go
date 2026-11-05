package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/sanitize"
	"portfolio/api/internal/storage"
)

var (
	hexColorRegex       = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
	ErrInvalidHexColor  = errors.New("brand_color harus berformat hex 6 karakter (contoh: #3178C6)")
	ErrInvalidGroupName = errors.New("group_name harus salah satu dari: language, framework, database, tool, cloud, other")
	ErrInvalidLevel     = errors.New("level harus bernilai antara 1 dan 5")
)

var allowedSkillGroups = map[string]bool{
	"language":  true,
	"framework": true,
	"database":  true,
	"tool":      true,
	"cloud":     true,
	"other":     true,
}

type SkillService interface {
	ListSkills(ctx context.Context, onlyVisible bool) ([]domain.Skill, error)
	GetSkillByID(ctx context.Context, id string) (*domain.Skill, error)
	CreateSkill(ctx context.Context, req domain.SkillCreateRequest) (*domain.Skill, error)
	UpdateSkill(ctx context.Context, id string, req domain.SkillUpdateRequest) (*domain.Skill, error)
	DeleteSkill(ctx context.Context, id string) error
	ValidateSkillIDs(ctx context.Context, skillIDs []string) ([]domain.Skill, error)
}

type skillService struct {
	repo repository.SkillRepository
}

func NewSkillService(repo repository.SkillRepository) SkillService {
	return &skillService{repo: repo}
}

func (s *skillService) ListSkills(ctx context.Context, onlyVisible bool) ([]domain.Skill, error) {
	return s.repo.List(ctx, onlyVisible)
}

func (s *skillService) GetSkillByID(ctx context.Context, id string) (*domain.Skill, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *skillService) CreateSkill(ctx context.Context, req domain.SkillCreateRequest) (*domain.Skill, error) {
	groupName := strings.ToLower(strings.TrimSpace(req.GroupName))
	if !allowedSkillGroups[groupName] {
		return nil, ErrInvalidGroupName
	}

	name := sanitize.NormalizeText(req.Name)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = GenerateSlug(name)
	} else {
		slug = GenerateSlug(slug)
	}

	// Validate Brand Color
	if req.BrandColor != nil && strings.TrimSpace(*req.BrandColor) != "" {
		trimmedColor := strings.TrimSpace(*req.BrandColor)
		if !hexColorRegex.MatchString(trimmedColor) {
			return nil, ErrInvalidHexColor
		}
		req.BrandColor = &trimmedColor
	} else {
		req.BrandColor = nil
	}

	// Validate Level
	if req.Level != nil && (*req.Level < 1 || *req.Level > 5) {
		return nil, ErrInvalidLevel
	}

	// Check name/slug conflict
	if existing, _ := s.repo.GetByName(ctx, name); existing != nil {
		return nil, repository.ErrSkillConflict
	}
	if existing, _ := s.repo.GetBySlug(ctx, slug); existing != nil {
		return nil, repository.ErrSkillConflict
	}

	isVisible := true
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}

	skill := &domain.Skill{
		Name:       name,
		Slug:       slug,
		GroupName:  groupName,
		IconKey:    storage.SanitizeStoragePathPtr(req.IconKey),
		BrandColor: req.BrandColor,
		Level:      req.Level,
		IsVisible:  isVisible,
		SortOrder:  req.SortOrder,
	}

	if err := s.repo.Create(ctx, skill); err != nil {
		return nil, fmt.Errorf("failed to create skill: %w", err)
	}

	return skill, nil
}

func (s *skillService) UpdateSkill(ctx context.Context, id string, req domain.SkillUpdateRequest) (*domain.Skill, error) {
	skill, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	groupName := strings.ToLower(strings.TrimSpace(req.GroupName))
	if !allowedSkillGroups[groupName] {
		return nil, ErrInvalidGroupName
	}

	name := sanitize.NormalizeText(req.Name)
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = GenerateSlug(name)
	} else {
		slug = GenerateSlug(slug)
	}

	if req.BrandColor != nil && strings.TrimSpace(*req.BrandColor) != "" {
		trimmedColor := strings.TrimSpace(*req.BrandColor)
		if !hexColorRegex.MatchString(trimmedColor) {
			return nil, ErrInvalidHexColor
		}
		req.BrandColor = &trimmedColor
	} else {
		req.BrandColor = nil
	}

	if req.Level != nil && (*req.Level < 1 || *req.Level > 5) {
		return nil, ErrInvalidLevel
	}

	// Check conflict if name or slug changed
	if name != skill.Name {
		if existing, _ := s.repo.GetByName(ctx, name); existing != nil && existing.ID != id {
			return nil, repository.ErrSkillConflict
		}
	}
	if slug != skill.Slug {
		if existing, _ := s.repo.GetBySlug(ctx, slug); existing != nil && existing.ID != id {
			return nil, repository.ErrSkillConflict
		}
	}

	isVisible := skill.IsVisible
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}

	skill.Name = name
	skill.Slug = slug
	skill.GroupName = groupName
	skill.IconKey = storage.SanitizeStoragePathPtr(req.IconKey)
	skill.BrandColor = req.BrandColor
	skill.Level = req.Level
	skill.IsVisible = isVisible
	skill.SortOrder = req.SortOrder

	if err := s.repo.Update(ctx, id, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *skillService) DeleteSkill(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *skillService) ValidateSkillIDs(ctx context.Context, skillIDs []string) ([]domain.Skill, error) {
	if len(skillIDs) == 0 {
		return []domain.Skill{}, nil
	}

	skills, err := s.repo.GetByIDs(ctx, skillIDs)
	if err != nil {
		return nil, err
	}

	if len(skills) != len(skillIDs) {
		return nil, errors.New("satu atau lebih skill_id tidak ditemukan")
	}

	return skills, nil
}
