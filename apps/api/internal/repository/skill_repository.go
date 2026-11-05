package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

var (
	ErrSkillNotFound = errors.New("skill tidak ditemukan")
	ErrSkillConflict = errors.New("skill dengan nama atau slug ini sudah terdaftar")
)

type SkillRepository interface {
	List(ctx context.Context, onlyVisible bool) ([]domain.Skill, error)
	GetByID(ctx context.Context, id string) (*domain.Skill, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Skill, error)
	GetByName(ctx context.Context, name string) (*domain.Skill, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.Skill, error)
	Create(ctx context.Context, s *domain.Skill) error
	Update(ctx context.Context, id string, s *domain.Skill) error
	Delete(ctx context.Context, id string) error
}

type postgresSkillRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSkillRepository(pool *pgxpool.Pool) SkillRepository {
	return &postgresSkillRepository{pool: pool}
}

func (r *postgresSkillRepository) List(ctx context.Context, onlyVisible bool) ([]domain.Skill, error) {
	var query string
	var rows pgx.Rows
	var err error

	if onlyVisible {
		query = `
			SELECT id, name, slug, group_name, icon_key, brand_color, level, is_visible, sort_order
			FROM skills
			WHERE is_visible = true
			ORDER BY group_name ASC, sort_order ASC, name ASC
		`
		rows, err = r.pool.Query(ctx, query)
	} else {
		query = `
			SELECT id, name, slug, group_name, icon_key, brand_color, level, is_visible, sort_order
			FROM skills
			ORDER BY group_name ASC, sort_order ASC, name ASC
		`
		rows, err = r.pool.Query(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []domain.Skill
	for rows.Next() {
		var s domain.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.GroupName, &s.IconKey, &s.BrandColor, &s.Level, &s.IsVisible, &s.SortOrder); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	if skills == nil {
		skills = []domain.Skill{}
	}
	return skills, nil
}

func (r *postgresSkillRepository) GetByID(ctx context.Context, id string) (*domain.Skill, error) {
	query := `
		SELECT id, name, slug, group_name, icon_key, brand_color, level, is_visible, sort_order
		FROM skills
		WHERE id = $1
	`
	var s domain.Skill
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.Slug, &s.GroupName, &s.IconKey, &s.BrandColor, &s.Level, &s.IsVisible, &s.SortOrder,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSkillNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresSkillRepository) GetBySlug(ctx context.Context, slug string) (*domain.Skill, error) {
	query := `
		SELECT id, name, slug, group_name, icon_key, brand_color, level, is_visible, sort_order
		FROM skills
		WHERE slug = $1
	`
	var s domain.Skill
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&s.ID, &s.Name, &s.Slug, &s.GroupName, &s.IconKey, &s.BrandColor, &s.Level, &s.IsVisible, &s.SortOrder,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSkillNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresSkillRepository) GetByName(ctx context.Context, name string) (*domain.Skill, error) {
	query := `
		SELECT id, name, slug, group_name, icon_key, brand_color, level, is_visible, sort_order
		FROM skills
		WHERE name = $1
	`
	var s domain.Skill
	err := r.pool.QueryRow(ctx, query, name).Scan(
		&s.ID, &s.Name, &s.Slug, &s.GroupName, &s.IconKey, &s.BrandColor, &s.Level, &s.IsVisible, &s.SortOrder,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSkillNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresSkillRepository) GetByIDs(ctx context.Context, ids []string) ([]domain.Skill, error) {
	if len(ids) == 0 {
		return []domain.Skill{}, nil
	}

	query := `
		SELECT id, name, slug, group_name, icon_key, brand_color, level, is_visible, sort_order
		FROM skills
		WHERE id = ANY($1)
	`
	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []domain.Skill
	for rows.Next() {
		var s domain.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.GroupName, &s.IconKey, &s.BrandColor, &s.Level, &s.IsVisible, &s.SortOrder); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	if skills == nil {
		skills = []domain.Skill{}
	}
	return skills, nil
}

func (r *postgresSkillRepository) Create(ctx context.Context, s *domain.Skill) error {
	query := `
		INSERT INTO skills (name, slug, group_name, icon_key, brand_color, level, is_visible, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return r.pool.QueryRow(ctx, query,
		s.Name, s.Slug, s.GroupName, s.IconKey, s.BrandColor, s.Level, s.IsVisible, s.SortOrder,
	).Scan(&s.ID)
}

func (r *postgresSkillRepository) Update(ctx context.Context, id string, s *domain.Skill) error {
	query := `
		UPDATE skills
		SET name = $1, slug = $2, group_name = $3, icon_key = $4, brand_color = $5,
		    level = $6, is_visible = $7, sort_order = $8
		WHERE id = $9
	`
	cmdTag, err := r.pool.Exec(ctx, query,
		s.Name, s.Slug, s.GroupName, s.IconKey, s.BrandColor, s.Level, s.IsVisible, s.SortOrder, id,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrSkillNotFound
	}
	return nil
}

func (r *postgresSkillRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM skills WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrSkillNotFound
	}
	return nil
}
