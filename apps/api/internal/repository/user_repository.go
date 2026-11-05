package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio/api/internal/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	RecordFailedLogin(ctx context.Context, id string, lockDuration time.Duration) error
	RecordSuccessfulLogin(ctx context.Context, id string) error
}

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{pool: pool}
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, failed_logins, locked_until, last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var u domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
		&u.FailedLogins, &u.LockedUntil, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, failed_logins, locked_until, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
		&u.FailedLogins, &u.LockedUntil, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresUserRepository) RecordFailedLogin(ctx context.Context, id string, lockDuration time.Duration) error {
	// Increment failed_logins. If >= 5, lock account for lockDuration.
	query := `
		UPDATE users
		SET failed_logins = failed_logins + 1,
		    locked_until = CASE 
		        WHEN failed_logins + 1 >= 5 THEN now() + $2::interval 
		        ELSE locked_until 
		    END,
		    updated_at = now()
		WHERE id = $1
	`
	intervalStr := lockDuration.String()
	_, err := r.pool.Exec(ctx, query, id, intervalStr)
	return err
}

func (r *postgresUserRepository) RecordSuccessfulLogin(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET failed_logins = 0,
		    locked_until = NULL,
		    last_login_at = now(),
		    updated_at = now()
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
