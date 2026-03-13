package handlers

import (
	"net/http"

	"booking-be/internal/observability"
	"booking-be/models"
	"booking-be/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SeatHandler struct {
	svc *service.SeatService
}

func NewSeatHandler(svc *service.SeatService) *SeatHandler {
	return &SeatHandler{svc: svc}
}

type generateSeatsRequest struct {
	Seats []models.Seat `json:"seats"`
}

func (h *SeatHandler) GenerateSeats(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	var req generateSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Str("trace_id", traceID).Str("event", "generate_seats_invalid_body").Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.GenerateSeats(c.Request.Context(), req.Seats); err != nil {
		log.Error().Str("trace_id", traceID).Str("event", "generate_seats_failed").Int("seat_count", len(req.Seats)).Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Str("trace_id", traceID).Str("event", "generate_seats_ok").Int("seat_count", len(req.Seats)).Send()
	c.JSON(http.StatusCreated, gin.H{"message": "seats saved", "count": len(req.Seats)})
}

func (h *SeatHandler) GetSeats(c *gin.Context) {
	traceID := observability.TraceIDFromContext(c.Request.Context())
	showtimeID := c.Param("showtimeId")
	if showtimeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "showtimeId is required"})
		return
	}
	seats, err := h.svc.GetSeats(c.Request.Context(), showtimeID)
	if err != nil {
		log.Error().Str("trace_id", traceID).Str("event", "get_seats_failed").Str("showtime_id", showtimeID).Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if seats == nil {
		seats = []models.Seat{}
	}
	log.Info().Str("trace_id", traceID).Str("event", "get_seats_ok").Str("showtime_id", showtimeID).Int("count", len(seats)).Send()
	c.JSON(http.StatusOK, gin.H{"seats": seats})
}
