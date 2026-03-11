package handlers

import (
	"net/http"

	"booking-be/service"

	"github.com/gin-gonic/gin"
)

// Handler depends on the service layer (DI)
type Handler struct {
	svc *service.Service
}

// NewHandler creates a new handler with the given service
func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
