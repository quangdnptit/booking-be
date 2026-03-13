package handlers

import (
	"net/http"
	"strings"

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
