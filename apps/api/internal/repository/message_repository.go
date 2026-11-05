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
	ErrMessageNotFound = errors.New("pesan tidak ditemukan")
)

type MessageRepository interface {
	List(ctx context.Context, status string, page, perPage int) ([]domain.Message, int, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
	Create(ctx context.Context, m *domain.Message) error
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
}

type postgresMessageRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresMessageRepository(pool *pgxpool.Pool) MessageRepository {
	return &postgresMessageRepository{pool: pool}
}

func (r *postgresMessageRepository) List(ctx context.Context, status string, page, perPage int) ([]domain.Message, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var conditions []string
	var args []interface{}
	argIdx := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(id) FROM messages %s`, whereClause)

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, subject, body, status, user_agent, created_at, read_at
		FROM messages
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, perPage, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(
			&m.ID, &m.Name, &m.Email, &m.Subject, &m.Body,
			&m.Status, &m.UserAgent, &m.CreatedAt, &m.ReadAt,
		); err != nil {
			return nil, 0, err
		}
		messages = append(messages, m)
	}

	if messages == nil {
		messages = []domain.Message{}
	}

	return messages, total, nil
}

func (r *postgresMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	query := `
		SELECT id, name, email, subject, body, status, user_agent, created_at, read_at
		FROM messages
		WHERE id = $1
	`
	var m domain.Message
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.Email, &m.Subject, &m.Body,
		&m.Status, &m.UserAgent, &m.CreatedAt, &m.ReadAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *postgresMessageRepository) Create(ctx context.Context, m *domain.Message) error {
	query := `
		INSERT INTO messages (name, email, subject, body, status, ip_hash, user_agent)
		VALUES ($1, $2, $3, $4, 'unread', $5, $6)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		m.Name, m.Email, m.Subject, m.Body, m.IPHash, m.UserAgent,
	).Scan(&m.ID, &m.CreatedAt)
}

func (r *postgresMessageRepository) UpdateStatus(ctx context.Context, id, status string) error {
	var query string
	if status == "read" {
		query = `
			UPDATE messages
			SET status = $1, read_at = CASE WHEN read_at IS NULL THEN now() ELSE read_at END
			WHERE id = $2
		`
	} else {
		query = `UPDATE messages SET status = $1 WHERE id = $2`
	}

	cmdTag, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrMessageNotFound
	}
	return nil
}

func (r *postgresMessageRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM messages WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrMessageNotFound
	}
	return nil
}
