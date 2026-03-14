package handlers

import (
	"errors"
	"net/http"

	"booking-be/internal/observability"
	"booking-be/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			log.Warn().Str("trace_id", traceID).Str("event", "auth_login_denied").Send()
			c.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrInvalidCredentials.Error()})
		case errors.Is(err, service.ErrAccountInactive):
			log.Warn().Str("trace_id", traceID).Str("event", "auth_login_inactive").Send()
			c.JSON(http.StatusForbidden, gin.H{"error": service.ErrAccountInactive.Error()})
		case errors.Is(err, service.ErrUserMisconfigured):
			log.Error().Str("trace_id", traceID).Str("event", "auth_login_misconfig").Send()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user record misconfigured"})
		default:
			if err != nil && err.Error() == "email and password required" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Error().Str("trace_id", traceID).Err(err).Str("event", "auth_login_failed").Send()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		}
		return
	}
	log.Info().Str("trace_id", traceID).Str("event", "auth_login_ok").Str("email", res.Email).Send()
	c.JSON(http.StatusOK, gin.H{
		"access_token": res.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   res.ExpiresIn,
		"user_id":      res.UserID,
		"email":        res.Email,
		"full_name":    res.FullName,
		"is_active":    res.IsActive,
		"amount":       res.Amount,
		"avatar":       res.Avatar,
		"created_at":   res.CreatedAt,
		"updated_at":   res.UpdatedAt,
	})
}

// Register POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.Register(c.Request.Context(), req.FullName, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyRegistered):
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		default:
			msg := err.Error()
			if msg == "full_name is required" || msg == "valid email is required" ||
				msg == "password must be at least 8 characters" {
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}
			log.Error().Str("trace_id", traceID).Err(err).Str("event", "auth_register_failed").Send()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register"})
		}
		return
	}
	log.Info().Str("trace_id", traceID).Str("event", "auth_register_ok").Str("email", res.Email).Send()
	c.JSON(http.StatusCreated, gin.H{
		"message":      "registered successfully",
		"email":        res.Email,
		"full_name":    res.FullName,
		"user_id":      res.UserID,
		"created_at":   res.CreatedAt,
		"updated_at":   res.UpdatedAt,
		"access_token": res.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   res.ExpiresIn,
	})
}
