package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/ratelimit"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/sanitize"
)

// ErrHoneypotTriggered is returned silently when the honeypot field is filled.
var ErrHoneypotTriggered = errors.New("bot terdeteksi")

// ErrContactRateLimited is returned when the IP sends too many contact requests.
var ErrContactRateLimited = errors.New("terlalu banyak pesan dikirim, coba lagi nanti")

// contactRateLimit: max 5 messages per hour per IP.
const contactRateLimit = 5
const contactRateWindow = time.Hour

type MessageService interface {
	CreateMessage(ctx context.Context, req domain.MessageCreateRequest, ip, userAgent string) (*domain.Message, error)
	ListMessages(ctx context.Context, status string, page, perPage int) ([]domain.Message, int, error)
	GetMessageByID(ctx context.Context, id string) (*domain.Message, error)
	UpdateMessageStatus(ctx context.Context, id, status string) error
	DeleteMessage(ctx context.Context, id string) error
}

type messageService struct {
	repo        repository.MessageRepository
	rateLimiter ratelimit.Limiter
	jwtSecret   string
}

func NewMessageService(repo repository.MessageRepository, rateLimiter ratelimit.Limiter, jwtSecret string) MessageService {
	return &messageService{
		repo:        repo,
		rateLimiter: rateLimiter,
		jwtSecret:   jwtSecret,
	}
}

// hashIP produces a deterministic, non-reversible identifier for an IP address.
// Uses SHA-256(ip + ":" + jwtSecret) so the mapping can't be reversed without the secret.
func (s *messageService) hashIP(ip string) string {
	h := sha256.Sum256([]byte(ip + ":" + s.jwtSecret))
	return fmt.Sprintf("%x", h)
}

func (s *messageService) CreateMessage(ctx context.Context, req domain.MessageCreateRequest, ip, userAgent string) (*domain.Message, error) {
	// 1. Honeypot check — silently reject bots that fill the hidden field.
	if strings.TrimSpace(req.HoneypotField) != "" {
		return nil, ErrHoneypotTriggered
	}

	// 2. Rate limit per IP (key: "contact:<ipHash>").
	ipHash := s.hashIP(ip)
	allowed, _ := s.rateLimiter.Allow("contact:"+ipHash, contactRateLimit, contactRateWindow)
	if !allowed {
		return nil, ErrContactRateLimited
	}

	// 3. Sanitize all fields — plain text only (no markdown).
	name := sanitize.NormalizeText(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	subject := sanitize.NormalizeText(req.Subject)
	body := sanitize.NormalizeText(req.Body)

	// 4. Re-validate body length after normalisation.
	bodyLen := utf8.RuneCountInString(body)
	if bodyLen < 10 {
		return nil, errors.New("isi pesan minimal 10 karakter")
	}
	if bodyLen > 5000 {
		return nil, errors.New("isi pesan maksimal 5000 karakter")
	}

	m := &domain.Message{
		Name:      name,
		Email:     email,
		Subject:   subject,
		Body:      body,
		IPHash:    ipHash,
		UserAgent: userAgent,
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("gagal menyimpan pesan: %w", err)
	}

	return m, nil
}

func (s *messageService) ListMessages(ctx context.Context, status string, page, perPage int) ([]domain.Message, int, error) {
	return s.repo.List(ctx, status, page, perPage)
}

func (s *messageService) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *messageService) UpdateMessageStatus(ctx context.Context, id, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *messageService) DeleteMessage(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
