package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// Context keys for handlers (after JWT middleware)
const (
	ContextUserID = "jwt_user_id"
)

var (
	ErrMissingSecret = errors.New("JWT_SECRET is not set")
	ErrMissingAuth   = errors.New("missing or invalid Authorization header")
	ErrInvalidToken  = errors.New("invalid or expired token")
)

// JWTAuthMiddleware validates Bearer JWT (HS256)
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	if strings.TrimSpace(secret) == "" {
		return func(c *gin.Context) {
			log.Error().Str("event", "jwt_middleware_misconfig").Msg("JWT_SECRET empty")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "server auth is not configured (JWT_SECRET)",
			})
		}
	}
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

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.Subject == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			return
		}

		c.Set(ContextUserID, claims.Subject)
		c.Next()
	}
}

// SignAccessToken issues a short-lived HS256 JWT with sub = userID.
func SignAccessToken(secret, userID string, ttl time.Duration) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrMissingSecret
	}
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
