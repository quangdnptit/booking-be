package repomodel

// BookingRecord is the persistence model for the DynamoDB bookings table.
// Used only within the repo layer for marshalling to/from DynamoDB.
type BookingRecord struct {
	ID         string `dynamodbav:"id"`
	UserID     string `dynamodbav:"user_id"`
	ShowtimeID string `dynamodbav:"showtime_id"`
	Status     string `dynamodbav:"status"`
	CreatedAt  string `dynamodbav:"created_at"`
}

// BookedSeatRecord is the persistence model for the DynamoDB booked_seats table.
// Used only within the repo layer for marshalling to/from DynamoDB.
type BookedSeatRecord struct {
	ID         string `dynamodbav:"id"`
	BookingID  string `dynamodbav:"booking_id"`
	ShowtimeID string `dynamodbav:"showtime_id"`
	SeatID     string `dynamodbav:"seat_id"`
	Status     string `dynamodbav:"status"`
}
