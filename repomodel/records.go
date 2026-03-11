package repomodel

// BookingRecord is the persistence model for the DynamoDB bookings table.
// Used only within the repo layer for marshalling to/from DynamoDB.
type BookingRecord struct {
	UserID      string  `dynamodbav:"user_id"`
	ShowtimeID  string  `dynamodbav:"showtime_id"`
	TotalAmount float64 `dynamodbav:"total_amount"`
	Status      string  `dynamodbav:"status"`
	CreatedAt   string  `dynamodbav:"created_at"`
	UpdatedAt   string  `dynamodbav:"updated_at"`
}

// BookedSeatRecord is the persistence model for the DynamoDB booked_seats table.
// Table key: pk=showtime_id, sk=seat_key.
type BookedSeatRecord struct {
	ShowtimeID string  `dynamodbav:"showtime_id"`
	SeatKey    string  `dynamodbav:"seat_key"`
	BookingID  string  `dynamodbav:"booking_id"`
	Status     string  `dynamodbav:"status"`
	Price      float32 `dynamodbav:"price"`
	CreatedAt  string  `dynamodbav:"created_at"`
	UpdatedAt  string  `dynamodbav:"updated_at"`
}
