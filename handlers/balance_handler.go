package handlers

import (
	"errors"
	"net/http"

	"booking-be/internal/auth"
	"booking-be/internal/observability"
	"booking-be/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type depositRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

// BalanceHandler handles deposit API.
type BalanceHandler struct {
	svc *service.BalanceService
}

func NewBalanceHandler(svc *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{svc: svc}
}

// Deposit POST /api/v1/deposit — add money to the authenticated user's balance.
func (h *BalanceHandler) Deposit(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	userID, _ := c.Get(auth.ContextUserID)
	uid, _ := userID.(string)

	var req depositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount is required and must be a number"})
		return
	}
	res, err := h.svc.Deposit(c.Request.Context(), uid, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDepositAmountInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrDepositAmountInvalid.Error()})
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			log.Error().Str("trace_id", traceID).Str("event", "deposit_failed").Err(err).Send()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "deposit failed"})
		}
		return
	}
	log.Info().Str("trace_id", traceID).Str("event", "deposit_ok").Str("user_id", uid).Float64("amount", req.Amount).Send()
	c.JSON(http.StatusOK, gin.H{
		"user_id":       res.UserID,
		"deposit_id":    res.DepositID,
		"amount":        res.Amount,
		"balance_after": res.BalanceAfter,
		"created_at":    res.CreatedAt,
	})
}
