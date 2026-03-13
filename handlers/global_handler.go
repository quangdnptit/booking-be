package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler depends on the service layer (DI)
type Handler struct {
}

// NewHandler creates a new handler with the given service
func NewHandler() *Handler {
	return &Handler{}
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
