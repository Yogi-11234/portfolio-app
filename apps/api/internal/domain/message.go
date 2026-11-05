package domain

import "time"

// Message represents an incoming contact form submission.
type Message struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Subject   string     `json:"subject"`
	Body      string     `json:"body"`
	Status    string     `json:"status"` // unread | read | archived | spam
	IPHash    string     `json:"-"`      // SHA-256(ip + secret), never exposed
	UserAgent string     `json:"user_agent,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

// MessageCreateRequest is the public contact form payload.
type MessageCreateRequest struct {
	Name    string `json:"name"    validate:"required,max=100"`
	Email   string `json:"email"   validate:"required,email,max=160"`
	Subject string `json:"subject" validate:"omitempty,max=160"`
	Body    string `json:"body"    validate:"required,min=10,max=5000"`
	// Honeypot field: bots tend to fill hidden inputs.
	// If non-empty the request is silently dropped.
	HoneypotField string `json:"website" validate:"-"`
}

// MessageStatusUpdateRequest allows admin to change a message's status.
type MessageStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=read archived spam unread"`
}
