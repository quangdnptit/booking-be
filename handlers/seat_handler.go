package handlers

import (
	"net/http"

	"booking-be/service"

	"github.com/gin-gonic/gin"
)

// SeatsGeneratorHandler depends on the service layer (DI)
type SeatsGeneratorHandler struct {
	svc *service.Service
}

// NewSeatsGeneratorHandler creates a new handler with the given service
func NewSeatsGeneratorHandler(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// SeatGenerate handles GET /SeatGenerate
func (h *Handler) SeatGenerate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
