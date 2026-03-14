package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"booking-be/internal/auth"
	"booking-be/repo"
	"booking-be/repomodel"

	"golang.org/x/crypto/bcrypt"
)

const minPasswordRunes = 8

var (
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrAccountInactive        = errors.New("account is not active")
	ErrUserMisconfigured      = errors.New("user record misconfigured")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

// LoginResult is what the handler serializes after a successful login.
type LoginResult struct {
	AccessToken string
	ExpiresIn   int
	UserID      string
	Email       string
	FullName    string
	IsActive    bool
	Amount      float64
	Avatar      string
	CreatedAt   string
	UpdatedAt   string
}

// RegisterResult is returned after successful registration.
type RegisterResult struct {
	UserID      string
	Email       string
	FullName    string
	CreatedAt   string
	UpdatedAt   string
	AccessToken string
	ExpiresIn   int
}

// AuthService orchestrates user auth (repo + bcrypt + JWT).
type AuthService struct {
	users  repo.UserRepo
	secret string
	ttl    time.Duration
}

func NewAuthService(users repo.UserRepo, jwtSecret string, ttl time.Duration) *AuthService {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &AuthService{users: users, secret: jwtSecret, ttl: ttl}
}

func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func looksLikeEmail(s string) bool {
	s = strings.TrimSpace(s)
	return strings.Contains(s, "@") && utf8.RuneCountInString(s) >= 5
}

func userIsActive(rec *repomodel.UserRecord) bool {
	if rec == nil {
		return false
	}
	s := strings.ToLower(strings.TrimSpace(rec.IsActive))
	return s == "" || s == "true" || s == "1" || s == "active"
}

// Login validates credentials and returns a JWT-backed result.
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password required")
	}

	rec, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("load user: %w", err)
	}
	if rec == nil || rec.PasswordHash == "" {
		return nil, ErrInvalidCredentials
	}
	if !userIsActive(rec) {
		return nil, ErrAccountInactive
	}
	if bcrypt.CompareHashAndPassword([]byte(rec.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	if rec.UserID == "" {
		return nil, ErrUserMisconfigured
	}

	now := time.Now().UTC().Format(time.RFC3339)
	_ = s.users.UpdateAudit(ctx, rec.Email, now) // best effort

	token, err := auth.SignAccessToken(s.secret, rec.UserID, s.ttl)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	return &LoginResult{
		AccessToken: token,
		ExpiresIn:   int(s.ttl.Seconds()),
		UserID:      rec.UserID,
		Email:       rec.Email,
		FullName:    rec.FullName,
		IsActive:    userIsActive(rec),
		Amount:      rec.Amount,
		Avatar:      rec.Avatar,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   now,
	}, nil
}

// Register creates a user with bcrypt password and returns a JWT
func (s *AuthService) Register(ctx context.Context, fullName, email, password string) (*RegisterResult, error) {
	fullName = strings.TrimSpace(fullName)
	email = normalizeEmail(email)
	if fullName == "" {
		return nil, fmt.Errorf("full_name is required")
	}
	if !looksLikeEmail(email) {
		return nil, fmt.Errorf("valid email is required")
	}
	if utf8.RuneCountInString(password) < minPasswordRunes {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	rec := repomodel.UserRecord{
		Email:        email,
		FullName:     fullName,
		PasswordHash: string(hash),
		UserID:       uuid.New().String(),
		IsActive:     "true",
		Amount:       0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(ctx, rec); err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailed") {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	token, err := auth.SignAccessToken(s.secret, rec.UserID, s.ttl)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}
	return &RegisterResult{
		UserID:      rec.UserID,
		Email:       rec.Email,
		FullName:    rec.FullName,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
		AccessToken: token,
		ExpiresIn:   int(s.ttl.Seconds()),
	}, nil
}
