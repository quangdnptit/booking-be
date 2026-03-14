package handlers

import (
	"net/http"
	"strings"

	"booking-be/internal/auth"
	"booking-be/internal/observability"
	"booking-be/models"
	"booking-be/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// BookingHandler handles booking HTTP API.
type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

// BookSeats handles POST /api/v1/bookings
func (h *BookingHandler) BookSeats(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	var req models.SeatsBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Str("trace_id", traceID).Str("event", "book_seats_invalid_body").Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sub, ok := c.Get(auth.ContextUserID); ok {
		if jwtUser, _ := sub.(string); jwtUser != "" && strings.TrimSpace(req.UserID) != jwtUser {
			c.JSON(http.StatusForbidden, gin.H{"error": "user_id must match authenticated user"})
			return
		}
	}

	booking, err := h.svc.BookSeats(c.Request.Context(), req)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "required"), strings.Contains(msg, "not found"), strings.Contains(msg, "not available"), strings.Contains(msg, "already held"):
			log.Warn().Str("trace_id", traceID).Str("event", "book_seats_bad_request").Err(err).Send()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		default:
			log.Error().Str("trace_id", traceID).Str("event", "book_seats_failed").Err(err).Send()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		return
	}

	log.Info().
		Str("trace_id", traceID).
		Str("event", "book_seats_ok").
		Str("booking_id", booking.ID).
		Int("seat_count", len(req.SeatKeys)).
		Send()
	c.JSON(http.StatusCreated, gin.H{
		"booking_id":   booking.ID,
		"user_id":      booking.UserID,
		"showtime_id":  booking.ShowtimeID,
		"total_amount": booking.TotalAmount,
		"status":       booking.Status,
		"seat_count":   len(req.SeatKeys),
	})
}

// GetUserBookingHistory handles GET /api/v1/users/:userId/bookings — list bookings from the bookings table.
func (h *BookingHandler) GetUserBookingHistory(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	userID := strings.TrimSpace(c.Param("userId"))
	if userID == "" {
		log.Warn().Str("trace_id", traceID).Str("event", "booking_history_missing_user").Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId path parameter is required"})
		return
	}
	if sub, ok := c.Get(auth.ContextUserID); ok {
		if jwtUser, _ := sub.(string); jwtUser != "" && userID != jwtUser {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot access another user's bookings"})
			return
		}
	}

	bookings, err := h.svc.GetUserBookingHistory(c.Request.Context(), userID)
	if err != nil {
		log.Error().Str("trace_id", traceID).Str("event", "booking_history_failed").Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Str("trace_id", traceID).Str("event", "booking_history_ok").Str("user_id", userID).Int("count", len(bookings)).Send()
	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"bookings": bookings,
		"count":    len(bookings),
	})
}
