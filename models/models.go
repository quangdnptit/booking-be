package models

// Room represents a hotel room
type Room struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Capacity      int     `json:"capacity"`
	PricePerNight float64 `json:"price_per_night"`
}

// Booking represents a room booking
type Booking struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status"`
}

// CreateBookingRequest represents the request payload for creating a booking
type CreateBookingRequest struct {
	RoomID    string `json:"room_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}
