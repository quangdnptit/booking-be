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
	ErrInvalidRefreshToken    = errors.New("invalid or expired refresh token")
)

// LoginResult is what the handler serializes after a successful login.
type LoginResult struct {
	AccessToken      string
	ExpiresIn        int
	RefreshToken     string
	RefreshExpiresIn int
	UserID           string
	Email            string
	FullName         string
	IsActive         bool
	Amount           float64
	Avatar           string
	CreatedAt        string
	UpdatedAt        string
}

// RegisterResult is returned after successful registration.
type RegisterResult struct {
	UserID           string
	Email            string
	FullName         string
	CreatedAt        string
	UpdatedAt        string
	AccessToken      string
	ExpiresIn        int
	RefreshToken     string
	RefreshExpiresIn int
}

// TokenPair is returned by Refresh.
type TokenPair struct {
	AccessToken      string
	ExpiresIn        int
	RefreshToken     string
	RefreshExpiresIn int
	TokenType        string
}

// AuthService orchestrates user auth (repo + bcrypt + JWT).
type AuthService struct {
	users      repo.UserRepo
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(users repo.UserRepo, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	if accessTTL <= 0 {
		accessTTL = time.Hour
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &AuthService{users: users, secret: jwtSecret, accessTTL: accessTTL, refreshTTL: refreshTTL}
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

func (s *AuthService) issueTokens(userID, email string) (access, refresh string, expIn, refreshExpIn int, err error) {
	access, err = auth.SignAccessToken(s.secret, userID, s.accessTTL)
	if err != nil {
		return "", "", 0, 0, err
	}
	refresh, err = auth.SignRefreshToken(s.secret, userID, email, s.refreshTTL)
	if err != nil {
		return "", "", 0, 0, err
	}
	return access, refresh, int(s.accessTTL.Seconds()), int(s.refreshTTL.Seconds()), nil
}

// Login validates credentials and returns access + refresh JWTs.
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
	_ = s.users.UpdateAudit(ctx, rec.Email, now)

	access, refresh, expIn, refreshExpIn, err := s.issueTokens(rec.UserID, rec.Email)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	return &LoginResult{
		AccessToken:      access,
		ExpiresIn:        expIn,
		RefreshToken:     refresh,
		RefreshExpiresIn: refreshExpIn,
		UserID:           rec.UserID,
		Email:            rec.Email,
		FullName:         rec.FullName,
		IsActive:         userIsActive(rec),
		Amount:           rec.Amount,
		Avatar:           rec.Avatar,
		CreatedAt:        rec.CreatedAt,
		UpdatedAt:        now,
	}, nil
}

// Register creates a user and returns access + refresh JWTs.
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

	access, refresh, expIn, refreshExpIn, err := s.issueTokens(rec.UserID, rec.Email)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}
	return &RegisterResult{
		UserID:           rec.UserID,
		Email:            rec.Email,
		FullName:         rec.FullName,
		CreatedAt:        rec.CreatedAt,
		UpdatedAt:        rec.UpdatedAt,
		AccessToken:      access,
		ExpiresIn:        expIn,
		RefreshToken:     refresh,
		RefreshExpiresIn: refreshExpIn,
	}, nil
}

// Refresh exchanges a valid refresh token for new access + refresh tokens (rotation).
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}
	userID, email, err := auth.ParseRefreshToken(s.secret, refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	rec, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("load user: %w", err)
	}
	if rec == nil || rec.UserID != userID {
		return nil, ErrInvalidRefreshToken
	}
	if !userIsActive(rec) {
		return nil, ErrAccountInactive
	}

	access, refresh, expIn, refreshExpIn, err := s.issueTokens(rec.UserID, rec.Email)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}
	return &TokenPair{
		AccessToken:      access,
		ExpiresIn:        expIn,
		RefreshToken:     refresh,
		RefreshExpiresIn: refreshExpIn,
		TokenType:        "Bearer",
	}, nil
}
