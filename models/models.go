package models

// Bookings is the domain model for a showtime booking
type Bookings struct {
	UserID      string  `dynamodbav:"user_id"`
	ShowtimeID  string  `dynamodbav:"showtime_id"`
	TotalAmount float64 `dynamodbav:"total_amount"`
	Status      string  `dynamodbav:"status"`
	CreatedAt   string  `dynamodbav:"created_at"`
	UpdatedAt   string  `dynamodbav:"updated_at"`
}

// BookedSeat is the domain model for a booked seat
type BookedSeat struct {
	ShowtimeID string  `dynamodbav:"showtime_id"`
	SeatKey    string  `dynamodbav:"seat_key"`
	BookingID  string  `dynamodbav:"booking_id"`
	Status     string  `dynamodbav:"status"`
	Price      float32 `dynamodbav:"price"`
	CreatedAt  string  `dynamodbav:"created_at"`
	UpdatedAt  string  `dynamodbav:"updated_at"`
}
