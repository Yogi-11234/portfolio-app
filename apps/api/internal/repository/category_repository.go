package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

var (
	ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
	ErrCategoryConflict = errors.New("kategori dengan tipe dan slug ini sudah terdaftar")
)

type CategoryRepository interface {
	List(ctx context.Context, categoryType string) ([]domain.Category, error)
	GetByID(ctx context.Context, id string) (*domain.Category, error)
	GetBySlugAndType(ctx context.Context, slug, categoryType string) (*domain.Category, error)
	Create(ctx context.Context, c *domain.Category) error
	Update(ctx context.Context, id string, c *domain.Category) error
	Delete(ctx context.Context, id string) error
}

type postgresCategoryRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &postgresCategoryRepository{pool: pool}
}

func (r *postgresCategoryRepository) List(ctx context.Context, categoryType string) ([]domain.Category, error) {
	var query string
	var rows pgx.Rows
	var err error

	if categoryType != "" {
		query = `
			SELECT id, name, slug, type, sort_order, created_at
			FROM categories
			WHERE type = $1
			ORDER BY sort_order ASC, name ASC
		`
		rows, err = r.pool.Query(ctx, query, categoryType)
	} else {
		query = `
			SELECT id, name, slug, type, sort_order, created_at
			FROM categories
			ORDER BY type ASC, sort_order ASC, name ASC
		`
		rows, err = r.pool.Query(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Type, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if categories == nil {
		categories = []domain.Category{}
	}
	return categories, nil
}

func (r *postgresCategoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `
		SELECT id, name, slug, type, sort_order, created_at
		FROM categories
		WHERE id = $1
	`
	var c domain.Category
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Type, &c.SortOrder, &c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *postgresCategoryRepository) GetBySlugAndType(ctx context.Context, slug, categoryType string) (*domain.Category, error) {
	query := `
		SELECT id, name, slug, type, sort_order, created_at
		FROM categories
		WHERE slug = $1 AND type = $2
	`
	var c domain.Category
	err := r.pool.QueryRow(ctx, query, slug, categoryType).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Type, &c.SortOrder, &c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *postgresCategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	query := `
		INSERT INTO categories (name, slug, type, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query, c.Name, c.Slug, c.Type, c.SortOrder).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresCategoryRepository) Update(ctx context.Context, id string, c *domain.Category) error {
	query := `
		UPDATE categories
		SET name = $1, slug = $2, sort_order = $3, type = $4
		WHERE id = $5
	`
	cmdTag, err := r.pool.Exec(ctx, query, c.Name, c.Slug, c.SortOrder, c.Type, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

func (r *postgresCategoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
