package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

var ErrProfileNotFound = errors.New("profile not found")

type ProfileRepository interface {
	Get(ctx context.Context) (*domain.Profile, error)
	GetByUserID(ctx context.Context, userID string) (*domain.Profile, error)
	Update(ctx context.Context, userID string, p *domain.Profile) error
}

type postgresProfileRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresProfileRepository(pool *pgxpool.Pool) ProfileRepository {
	return &postgresProfileRepository{pool: pool}
}

func (r *postgresProfileRepository) Get(ctx context.Context) (*domain.Profile, error) {
	query := `
		SELECT id, user_id, full_name, headline, bio_md, bio_html, location, public_email,
		       avatar_path, resume_path, github_url, linkedin_url, social_links, open_to_work, updated_at
		FROM profile
		LIMIT 1
	`
	var p domain.Profile
	err := r.pool.QueryRow(ctx, query).Scan(
		&p.ID, &p.UserID, &p.FullName, &p.Headline, &p.BioMD, &p.BioHTML, &p.Location, &p.PublicEmail,
		&p.AvatarPath, &p.ResumePath, &p.GithubURL, &p.LinkedinURL, &p.SocialLinks, &p.OpenToWork, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresProfileRepository) GetByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	query := `
		SELECT id, user_id, full_name, headline, bio_md, bio_html, location, public_email,
		       avatar_path, resume_path, github_url, linkedin_url, social_links, open_to_work, updated_at
		FROM profile
		WHERE user_id = $1
	`
	var p domain.Profile
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.FullName, &p.Headline, &p.BioMD, &p.BioHTML, &p.Location, &p.PublicEmail,
		&p.AvatarPath, &p.ResumePath, &p.GithubURL, &p.LinkedinURL, &p.SocialLinks, &p.OpenToWork, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresProfileRepository) Update(ctx context.Context, userID string, p *domain.Profile) error {
	query := `
		UPDATE profile
		SET full_name = $1, headline = $2, bio_md = $3, bio_html = $4, location = $5,
		    public_email = $6, avatar_path = $7, resume_path = $8, github_url = $9,
		    linkedin_url = $10, social_links = $11, open_to_work = $12, updated_at = now()
		WHERE user_id = $13
		RETURNING id, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		p.FullName, p.Headline, p.BioMD, p.BioHTML, p.Location,
		p.PublicEmail, p.AvatarPath, p.ResumePath, p.GithubURL,
		p.LinkedinURL, p.SocialLinks, p.OpenToWork, userID,
	).Scan(&p.ID, &p.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProfileNotFound
		}
		return err
	}
	return nil
}
