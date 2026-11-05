package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"portfolio/api/internal/config"
	"portfolio/api/internal/domain"
	"portfolio/api/internal/repository"
	"portfolio/api/internal/security"
)

var (
	ErrInvalidCredentials = errors.New("email atau password tidak valid")
	ErrAccountLocked      = errors.New("akun terkunci sementara karena terlalu banyak percobaan gagal. Silakan coba 15 menit lagi")
	ErrInvalidToken       = errors.New("token otorisasi tidak valid")
	ErrTokenExpired       = errors.New("sesi telah kedaluwarsa, silakan login kembali")
)

type JWTClaims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(ctx context.Context, req domain.LoginRequest) (*domain.User, string, string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GenerateCSRFToken() (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.User, string, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", err
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, "", "", ErrAccountLocked
	}

	// Verify password using Argon2id
	match, err := security.ComparePassword(req.Password, user.PasswordHash)
	if err != nil || !match {
		// Record failed attempt (5 failures lock account for 15 minutes)
		_ = s.userRepo.RecordFailedLogin(ctx, user.ID, 15*time.Minute)
		return nil, "", "", ErrInvalidCredentials
	}

	// Password is correct, reset failed logins
	if err := s.userRepo.RecordSuccessfulLogin(ctx, user.ID); err != nil {
		return nil, "", "", fmt.Errorf("failed to update login stats: %w", err)
	}

	// Generate JWT
	ttl := time.Duration(s.cfg.JWTTLHours) * time.Hour
	claims := JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "portfolio-api",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	// Generate CSRF Token for double-submit
	csrfToken, err := s.GenerateCSRFToken()
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate csrf token: %w", err)
	}

	return user, tokenString, csrfToken, nil
}

func (s *authService) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Strict algorithm check
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if token.Header["alg"] != "HS256" {
			return nil, fmt.Errorf("algorithm %v is rejected, only HS256 allowed", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *authService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *authService) GenerateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
