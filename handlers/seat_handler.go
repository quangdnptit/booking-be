package handlers

import (
	"net/http"

	"booking-be/models"
	"booking-be/service"

	"github.com/gin-gonic/gin"
)

// SeatHandler depends on the service layer (DI)
type SeatHandler struct {
	svc *service.SeatService
}

// NewSeatHandler creates a new handler with the given service
func NewSeatHandler(svc *service.SeatService) *SeatHandler {
	return &SeatHandler{
		svc: svc,
	}
}

// generateSeatsRequest is the JSON body for POST /api/v1/seats/generate-seats
type generateSeatsRequest struct {
	Seats []models.Seat `json:"seats"`
}

// SeatGenerate reads seats from the request body and batch-saves them.
// Body example:
//
//	{
//	  "seats": [
//	    {
//	      "showtime_id": "st-1",
//	      "seat_key": "A#1",
//	      "room_id": "550e8400-e29b-41d4-a716-446655440000",
//	      "seat_type": "STANDARD",
//	      "booking_id": "",
//	      "is_active": "true",
//	      "price": 12.5,
//	      "seat_status": "AVAILABLE"
//	    }
//	  ]
//	}
func (h *SeatHandler) GenerateSeats(c *gin.Context) {
	var req generateSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.GenerateSeats(c.Request.Context(), req.Seats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "seats saved", "count": len(req.Seats)})
}
