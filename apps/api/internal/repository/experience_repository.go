package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

var ErrExperienceNotFound = errors.New("pengalaman tidak ditemukan")

type ExperienceRepository interface {
	List(ctx context.Context, onlyVisible bool, expType string) ([]domain.Experience, error)
	GetByID(ctx context.Context, id string) (*domain.Experience, error)
	Create(ctx context.Context, e *domain.Experience) error
	Update(ctx context.Context, id string, e *domain.Experience) error
	Delete(ctx context.Context, id string) error
}

type postgresExperienceRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresExperienceRepository(pool *pgxpool.Pool) ExperienceRepository {
	return &postgresExperienceRepository{pool: pool}
}

func (r *postgresExperienceRepository) List(ctx context.Context, onlyVisible bool, expType string) ([]domain.Experience, error) {
	query := `
		SELECT id, type, organization, role, location, description_md, description_html,
		       start_date, end_date, is_current, is_visible, sort_order, created_at, updated_at
		FROM experiences
		WHERE ($1 = false OR is_visible = true)
		  AND ($2 = '' OR type = $2)
		ORDER BY sort_order ASC, start_date DESC
	`
	rows, err := r.pool.Query(ctx, query, onlyVisible, expType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Experience
	for rows.Next() {
		var e domain.Experience
		err := rows.Scan(
			&e.ID, &e.Type, &e.Organization, &e.Role, &e.Location, &e.DescriptionMD, &e.DescriptionHTML,
			&e.StartDate, &e.EndDate, &e.IsCurrent, &e.IsVisible, &e.SortOrder, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}

	if list == nil {
		list = []domain.Experience{}
	}
	return list, nil
}

func (r *postgresExperienceRepository) GetByID(ctx context.Context, id string) (*domain.Experience, error) {
	query := `
		SELECT id, type, organization, role, location, description_md, description_html,
		       start_date, end_date, is_current, is_visible, sort_order, created_at, updated_at
		FROM experiences
		WHERE id = $1
	`
	var e domain.Experience
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.Type, &e.Organization, &e.Role, &e.Location, &e.DescriptionMD, &e.DescriptionHTML,
		&e.StartDate, &e.EndDate, &e.IsCurrent, &e.IsVisible, &e.SortOrder, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExperienceNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *postgresExperienceRepository) Create(ctx context.Context, e *domain.Experience) error {
	query := `
		INSERT INTO experiences (type, organization, role, location, description_md, description_html,
		                        start_date, end_date, is_current, is_visible, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		e.Type, e.Organization, e.Role, e.Location, e.DescriptionMD, e.DescriptionHTML,
		e.StartDate, e.EndDate, e.IsCurrent, e.IsVisible, e.SortOrder,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *postgresExperienceRepository) Update(ctx context.Context, id string, e *domain.Experience) error {
	query := `
		UPDATE experiences
		SET type = $1, organization = $2, role = $3, location = $4, description_md = $5,
		    description_html = $6, start_date = $7, end_date = $8, is_current = $9,
		    is_visible = $10, sort_order = $11, updated_at = now()
		WHERE id = $12
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		e.Type, e.Organization, e.Role, e.Location, e.DescriptionMD, e.DescriptionHTML,
		e.StartDate, e.EndDate, e.IsCurrent, e.IsVisible, e.SortOrder, id,
	).Scan(&e.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrExperienceNotFound
		}
		return err
	}
	return nil
}

func (r *postgresExperienceRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM experiences WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrExperienceNotFound
	}
	return nil
}
