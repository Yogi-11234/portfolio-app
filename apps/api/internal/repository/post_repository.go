package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

var (
	ErrPostNotFound = errors.New("postingan blog tidak ditemukan")
	ErrPostConflict = errors.New("postingan dengan slug ini sudah terdaftar")
)

type PostRepository interface {
	List(ctx context.Context, f domain.PostFilter) ([]domain.Post, int, error)
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Post, error)
	Create(ctx context.Context, p *domain.Post) error
	Update(ctx context.Context, id string, p *domain.Post) error
	Delete(ctx context.Context, id string) error
}

type postgresPostRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPostRepository(pool *pgxpool.Pool) PostRepository {
	return &postgresPostRepository{pool: pool}
}

func (r *postgresPostRepository) List(ctx context.Context, f domain.PostFilter) ([]domain.Post, int, error) {
	page := f.Page
	if page < 1 {
		page = 1
	}
	perPage := f.PerPage
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	var conditions []string
	var args []interface{}
	argIdx := 1

	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argIdx))
		args = append(args, f.Status)
		argIdx++
	}

	if f.CategorySlug != "" {
		conditions = append(conditions, fmt.Sprintf("c.slug = $%d", argIdx))
		args = append(args, f.CategorySlug)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(p.id)
		FROM posts p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
	`, whereClause)

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.author_id, p.category_id, p.title, p.slug, p.excerpt, p.content_md, p.content_html,
		       p.cover_path, p.status, p.reading_minutes, p.published_at, p.created_at, p.updated_at,
		       c.id, c.name, c.slug,
		       u.name
		FROM posts p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN users u ON p.author_id = u.id
		%s
		ORDER BY COALESCE(p.published_at, p.created_at) DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, perPage, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var p domain.Post
		var catID, catName, catSlug *string

		err := rows.Scan(
			&p.ID, &p.AuthorID, &p.CategoryID, &p.Title, &p.Slug, &p.Excerpt, &p.ContentMD, &p.ContentHTML,
			&p.CoverPath, &p.Status, &p.ReadingMinutes, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
			&catID, &catName, &catSlug,
			&p.AuthorName,
		)
		if err != nil {
			return nil, 0, err
		}

		if catID != nil && catName != nil && catSlug != nil {
			p.Category = &domain.CategoryBrief{
				ID:   *catID,
				Name: *catName,
				Slug: *catSlug,
			}
		}

		posts = append(posts, p)
	}

	if posts == nil {
		posts = []domain.Post{}
	}

	return posts, total, nil
}

func (r *postgresPostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	query := `
		SELECT p.id, p.author_id, p.category_id, p.title, p.slug, p.excerpt, p.content_md, p.content_html,
		       p.cover_path, p.status, p.reading_minutes, p.published_at, p.created_at, p.updated_at,
		       c.id, c.name, c.slug,
		       u.name
		FROM posts p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.id = $1
	`
	var p domain.Post
	var catID, catName, catSlug *string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.AuthorID, &p.CategoryID, &p.Title, &p.Slug, &p.Excerpt, &p.ContentMD, &p.ContentHTML,
		&p.CoverPath, &p.Status, &p.ReadingMinutes, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
		&catID, &catName, &catSlug,
		&p.AuthorName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	if catID != nil && catName != nil && catSlug != nil {
		p.Category = &domain.CategoryBrief{
			ID:   *catID,
			Name: *catName,
			Slug: *catSlug,
		}
	}

	return &p, nil
}

func (r *postgresPostRepository) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	query := `
		SELECT p.id, p.author_id, p.category_id, p.title, p.slug, p.excerpt, p.content_md, p.content_html,
		       p.cover_path, p.status, p.reading_minutes, p.published_at, p.created_at, p.updated_at,
		       c.id, c.name, c.slug,
		       u.name
		FROM posts p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.slug = $1
	`
	var p domain.Post
	var catID, catName, catSlug *string

	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.AuthorID, &p.CategoryID, &p.Title, &p.Slug, &p.Excerpt, &p.ContentMD, &p.ContentHTML,
		&p.CoverPath, &p.Status, &p.ReadingMinutes, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
		&catID, &catName, &catSlug,
		&p.AuthorName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	if catID != nil && catName != nil && catSlug != nil {
		p.Category = &domain.CategoryBrief{
			ID:   *catID,
			Name: *catName,
			Slug: *catSlug,
		}
	}

	return &p, nil
}

func (r *postgresPostRepository) Create(ctx context.Context, p *domain.Post) error {
	query := `
		INSERT INTO posts (author_id, category_id, title, slug, excerpt, content_md, content_html,
		                  cover_path, status, reading_minutes, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		p.AuthorID, p.CategoryID, p.Title, p.Slug, p.Excerpt, p.ContentMD, p.ContentHTML,
		p.CoverPath, p.Status, p.ReadingMinutes, p.PublishedAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *postgresPostRepository) Update(ctx context.Context, id string, p *domain.Post) error {
	query := `
		UPDATE posts
		SET category_id = $1, title = $2, slug = $3, excerpt = $4, content_md = $5,
		    content_html = $6, cover_path = $7, status = $8, reading_minutes = $9,
		    published_at = $10, updated_at = now()
		WHERE id = $11
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		p.CategoryID, p.Title, p.Slug, p.Excerpt, p.ContentMD,
		p.ContentHTML, p.CoverPath, p.Status, p.ReadingMinutes,
		p.PublishedAt, id,
	).Scan(&p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPostNotFound
		}
		return err
	}
	return nil
}

func (r *postgresPostRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrPostNotFound
	}
	return nil
}
