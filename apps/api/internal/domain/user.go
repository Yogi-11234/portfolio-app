package domain

import "time"

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Name         string     `json:"name"`
	Role         string     `json:"role"`
	FailedLogins int16      `json:"-"`
	LockedUntil  *time.Time `json:"-"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=160"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type UserResponse struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Role        string     `json:"role"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Role:        u.Role,
		LastLoginAt: u.LastLoginAt,
	}
}
