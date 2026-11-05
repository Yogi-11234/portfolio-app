package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

var ErrCertificateNotFound = errors.New("sertifikat tidak ditemukan")

type CertificateRepository interface {
	List(ctx context.Context, onlyVisible bool) ([]domain.Certificate, error)
	GetByID(ctx context.Context, id string) (*domain.Certificate, error)
	Create(ctx context.Context, c *domain.Certificate) error
	Update(ctx context.Context, id string, c *domain.Certificate) error
	Delete(ctx context.Context, id string) error
}

type postgresCertificateRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresCertificateRepository(pool *pgxpool.Pool) CertificateRepository {
	return &postgresCertificateRepository{pool: pool}
}

func (r *postgresCertificateRepository) List(ctx context.Context, onlyVisible bool) ([]domain.Certificate, error) {
	query := `
		SELECT id, title, issuer, credential_id, credential_url, issued_at, expires_at,
		       file_path, file_mime, thumbnail_path, is_visible, sort_order, created_at, updated_at
		FROM certificates
		WHERE ($1 = false OR is_visible = true)
		ORDER BY sort_order ASC, issued_at DESC
	`
	rows, err := r.pool.Query(ctx, query, onlyVisible)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Certificate
	for rows.Next() {
		var c domain.Certificate
		err := rows.Scan(
			&c.ID, &c.Title, &c.Issuer, &c.CredentialID, &c.CredentialURL,
			&c.IssuedAt, &c.ExpiresAt, &c.FilePath, &c.FileMime, &c.ThumbnailPath,
			&c.IsVisible, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	if list == nil {
		list = []domain.Certificate{}
	}
	return list, nil
}

func (r *postgresCertificateRepository) GetByID(ctx context.Context, id string) (*domain.Certificate, error) {
	query := `
		SELECT id, title, issuer, credential_id, credential_url, issued_at, expires_at,
		       file_path, file_mime, thumbnail_path, is_visible, sort_order, created_at, updated_at
		FROM certificates
		WHERE id = $1
	`
	var c domain.Certificate
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Title, &c.Issuer, &c.CredentialID, &c.CredentialURL,
		&c.IssuedAt, &c.ExpiresAt, &c.FilePath, &c.FileMime, &c.ThumbnailPath,
		&c.IsVisible, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCertificateNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *postgresCertificateRepository) Create(ctx context.Context, c *domain.Certificate) error {
	query := `
		INSERT INTO certificates (title, issuer, credential_id, credential_url, issued_at, expires_at,
		                         file_path, file_mime, thumbnail_path, is_visible, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		c.Title, c.Issuer, c.CredentialID, c.CredentialURL, c.IssuedAt, c.ExpiresAt,
		c.FilePath, c.FileMime, c.ThumbnailPath, c.IsVisible, c.SortOrder,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *postgresCertificateRepository) Update(ctx context.Context, id string, c *domain.Certificate) error {
	query := `
		UPDATE certificates
		SET title = $1, issuer = $2, credential_id = $3, credential_url = $4, issued_at = $5,
		    expires_at = $6, file_path = $7, file_mime = $8, thumbnail_path = $9,
		    is_visible = $10, sort_order = $11, updated_at = now()
		WHERE id = $12
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		c.Title, c.Issuer, c.CredentialID, c.CredentialURL, c.IssuedAt, c.ExpiresAt,
		c.FilePath, c.FileMime, c.ThumbnailPath, c.IsVisible, c.SortOrder, id,
	).Scan(&c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCertificateNotFound
		}
		return err
	}
	return nil
}

func (r *postgresCertificateRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM certificates WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrCertificateNotFound
	}
	return nil
}
