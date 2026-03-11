package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"booking-be/models"
	"booking-be/storage"
)

// Handler contains the storage and provides API handlers
type Handler struct {
	store storage.Store
}

// NewHandler creates a new handler with the given store
func NewHandler(store storage.Store) *Handler {
	return &Handler{
		store: store,
	}
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

// ListRooms handles GET /rooms
func (h *Handler) ListRooms(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.GetRooms())
}

// GetRoom handles GET /rooms/:id
func (h *Handler) GetRoom(c *gin.Context) {
	roomID := c.Param("id")
	room, found := h.store.GetRoomByID(roomID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// CheckAvailability handles GET /rooms/:id/availability
func (h *Handler) CheckAvailability(c *gin.Context) {
	roomID := c.Param("id")
	_, found := h.store.GetRoomByID(roomID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"room_id":      roomID,
		"is_available": true,
	})
}

// CreateBooking handles POST /bookings
func (h *Handler) CreateBooking(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking := models.Booking{
		ID:        uuid.New().String(),
		RoomID:    req.RoomID,
		UserID:    req.UserID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    "CONFIRMED",
	}

	result := h.store.CreateBooking(booking)
	c.JSON(http.StatusCreated, result)
}

// ListBookings handles GET /bookings
func (h *Handler) ListBookings(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.GetBookings())
}

// GetBooking handles GET /bookings/:id
func (h *Handler) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	booking, found := h.store.GetBookingByID(bookingID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}
	c.JSON(http.StatusOK, booking)
}
