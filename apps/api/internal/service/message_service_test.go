package service

import (
	"context"
	"errors"
	"testing"

	"portfolio/api/internal/domain"
	"portfolio/api/internal/ratelimit"
	"portfolio/api/internal/repository"
)

// --- Mock repository ---

type mockMessageRepo struct {
	messages map[string]*domain.Message
	nextID   int
}

func newMockMessageRepo() *mockMessageRepo {
	return &mockMessageRepo{messages: make(map[string]*domain.Message)}
}

func (m *mockMessageRepo) List(_ context.Context, _ string, _, _ int) ([]domain.Message, int, error) {
	return nil, 0, nil
}
func (m *mockMessageRepo) GetByID(_ context.Context, id string) (*domain.Message, error) {
	msg, ok := m.messages[id]
	if !ok {
		return nil, repository.ErrMessageNotFound
	}
	return msg, nil
}
func (m *mockMessageRepo) Create(_ context.Context, msg *domain.Message) error {
	m.nextID++
	msg.ID = "msg-" + string(rune('0'+m.nextID))
	m.messages[msg.ID] = msg
	return nil
}
func (m *mockMessageRepo) UpdateStatus(_ context.Context, id, status string) error {
	msg, ok := m.messages[id]
	if !ok {
		return repository.ErrMessageNotFound
	}
	msg.Status = status
	return nil
}
func (m *mockMessageRepo) Delete(_ context.Context, id string) error {
	if _, ok := m.messages[id]; !ok {
		return repository.ErrMessageNotFound
	}
	delete(m.messages, id)
	return nil
}

// --- Helpers ---

func newTestMessageService(t *testing.T) (MessageService, *mockMessageRepo) {
	t.Helper()
	repo := newMockMessageRepo()
	rl := ratelimit.NewMemoryLimiter()
	svc := NewMessageService(repo, rl, "test-secret-key-for-hashing-only")
	return svc, repo
}

func validMsgReq() domain.MessageCreateRequest {
	return domain.MessageCreateRequest{
		Name:  "Budi Santoso",
		Email: "budi@example.com",
		Body:  "Halo, saya tertarik dengan proyek Anda dan ingin berdiskusi lebih lanjut.",
	}
}

// --- Tests ---

func TestMessageService_HoneypotRejectsBot(t *testing.T) {
	svc, _ := newTestMessageService(t)
	req := validMsgReq()
	req.HoneypotField = "http://spam.example.com" // Bot filled the hidden field

	_, err := svc.CreateMessage(context.Background(), req, "1.2.3.4", "bot/1.0")
	if !errors.Is(err, ErrHoneypotTriggered) {
		t.Fatalf("expected ErrHoneypotTriggered, got: %v", err)
	}
}

func TestMessageService_RateLimitContact(t *testing.T) {
	svc, _ := newTestMessageService(t)
	ip := "10.0.0.1"

	// Send contactRateLimit (5) messages successfully.
	for i := 0; i < contactRateLimit; i++ {
		req := validMsgReq()
		req.Email = "user@example.com"
		if _, err := svc.CreateMessage(context.Background(), req, ip, "UA"); err != nil {
			t.Fatalf("unexpected error on attempt %d: %v", i+1, err)
		}
	}

	// 6th attempt should be rate-limited.
	req := validMsgReq()
	_, err := svc.CreateMessage(context.Background(), req, ip, "UA")
	if !errors.Is(err, ErrContactRateLimited) {
		t.Fatalf("expected ErrContactRateLimited, got: %v", err)
	}
}

func TestMessageService_ValidMessage(t *testing.T) {
	svc, repo := newTestMessageService(t)
	req := validMsgReq()
	req.Subject = "Kolaborasi proyek"

	msg, err := svc.CreateMessage(context.Background(), req, "192.168.1.1", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.ID == "" {
		t.Fatal("expected non-empty message ID")
	}
	if msg.IPHash == "" {
		t.Fatal("expected non-empty ip_hash")
	}
	// IPHash should NOT be the raw IP.
	if msg.IPHash == "192.168.1.1" {
		t.Fatal("ip_hash must not be the raw IP address")
	}
	// Message should be persisted in the mock repo.
	stored, err := repo.GetByID(context.Background(), msg.ID)
	if err != nil {
		t.Fatalf("message not found in repo: %v", err)
	}
	if stored.Name != "Budi Santoso" {
		t.Fatalf("unexpected name: %s", stored.Name)
	}
}

func TestMessageService_BodyTooShort(t *testing.T) {
	svc, _ := newTestMessageService(t)
	req := validMsgReq()
	req.Body = "Hi" // only 2 runes

	_, err := svc.CreateMessage(context.Background(), req, "1.1.1.1", "UA")
	if err == nil {
		t.Fatal("expected error for body < 10 chars, got nil")
	}
}

func TestMessageService_UpdateStatus(t *testing.T) {
	svc, repo := newTestMessageService(t)
	// Seed a message directly in the mock repo.
	seed := &domain.Message{ID: "msg-seed", Status: "unread"}
	repo.messages["msg-seed"] = seed

	if err := svc.UpdateMessageStatus(context.Background(), "msg-seed", "read"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.messages["msg-seed"].Status != "read" {
		t.Fatal("expected status to be 'read'")
	}
}

func TestMessageService_DeleteMessage(t *testing.T) {
	svc, repo := newTestMessageService(t)
	repo.messages["msg-del"] = &domain.Message{ID: "msg-del"}

	if err := svc.DeleteMessage(context.Background(), "msg-del"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := repo.messages["msg-del"]; exists {
		t.Fatal("expected message to be deleted")
	}

	// Deleting again should return ErrMessageNotFound.
	if err := svc.DeleteMessage(context.Background(), "msg-del"); !errors.Is(err, repository.ErrMessageNotFound) {
		t.Fatalf("expected ErrMessageNotFound, got: %v", err)
	}
}
