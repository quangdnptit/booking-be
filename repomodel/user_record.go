package repomodel

// UserRecord is stored in DynamoDB table "users". Partition key: email (identity).
type UserRecord struct {
	Email        string  `dynamo:"email"`
	FullName     string  `dynamo:"full_name"` // display only
	PasswordHash string  `dynamo:"password_hash"`
	UserID       string  `dynamo:"user_id"` // JWT sub; stable id for bookings
	IsActive     string  `dynamo:"is_active"`
	Amount       float64 `dynamo:"amount"`
	Avatar       string  `dynamo:"avatar"`
	CreatedAt    string  `dynamo:"created_at"`
	UpdatedAt    string  `dynamo:"updated_at"`
}
