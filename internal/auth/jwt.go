package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	claimType   = "typ"
	typeAccess  = "access"
	typeRefresh = "refresh"
	claimEmail  = "email"
)

// Context keys for handlers (after JWT middleware)
const (
	ContextUserID = "jwt_user_id"
)

var (
	ErrMissingAuth     = errors.New("missing or invalid Authorization header")
	ErrInvalidToken    = errors.New("invalid or expired token")
	ErrInvalidTokenUse = errors.New("wrong token type; use access token")
)

// JWTAuthMiddleware validates Bearer access JWT only (HS256, typ=access).
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrMissingAuth.Error()})
			return
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(raw, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrMissingAuth.Error()})
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(raw, prefix))
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrMissingAuth.Error()})
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			return
		}
		if typ, _ := claims[claimType].(string); typ != typeAccess {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidTokenUse.Error()})
			return
		}
		sub, _ := claims["sub"].(string)
		if sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			return
		}

		c.Set(ContextUserID, sub)
		c.Next()
	}
}

// SignAccessToken issues a short-lived HS256 JWT (typ=access).
func SignAccessToken(secret, userID string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":     userID,
		claimType: typeAccess,
		"iat":     jwt.NewNumericDate(now).Unix(),
		"exp":     jwt.NewNumericDate(now.Add(ttl)).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// SignRefreshToken issues a longer-lived HS256 JWT (typ=refresh); email used to reload user on refresh.
func SignRefreshToken(secret, userID, email string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      userID,
		claimType:  typeRefresh,
		claimEmail: strings.ToLower(strings.TrimSpace(email)),
		"iat":      jwt.NewNumericDate(now).Unix(),
		"exp":      jwt.NewNumericDate(now.Add(ttl)).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// ParseRefreshToken validates a refresh JWT and returns userID and email.
func ParseRefreshToken(secret, tokenStr string) (userID, email string, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", "", ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", ErrInvalidToken
	}
	if typ, _ := claims[claimType].(string); typ != typeRefresh {
		return "", "", ErrInvalidToken
	}
	userID, _ = claims["sub"].(string)
	email, _ = claims[claimEmail].(string)
	if userID == "" || email == "" {
		return "", "", ErrInvalidToken
	}
	return userID, email, nil
}
