package repomodel

// Single-table design: user profile and deposit rows share the same table
// Table keys: pk (string), sk (string)
// GSI "email-index": partition key = email, sort key = sk (for GetByEmail)

const (
	PKPrefixUser    = "USER#"
	PKPrefixBalance = "BALANCE#"
	SKProfile       = "PROFILE"
	SKDepositPrefix = "DEPOSIT#" // sk=DEPOSIT#<deposit_id> so all deposits are kept for history
)

// UserRecord is the user profile row. pk=USER#<user_id>, sk=PROFILE
type UserRecord struct {
	Pk           string  `dynamo:"pk"`
	Sk           string  `dynamo:"sk"`
	Email        string  `dynamo:"email"`
	FullName     string  `dynamo:"full_name"`
	PasswordHash string  `dynamo:"password_hash"`
	UserID       string  `dynamo:"user_id"`
	IsActive     string  `dynamo:"is_active"`
	Amount       float64 `dynamo:"amount"`
	Avatar       string  `dynamo:"avatar"`
	CreatedAt    string  `dynamo:"created_at"`
	UpdatedAt    string  `dynamo:"updated_at"`
}

// BalanceRecord is one deposit row; all are kept for history. pk=BALANCE#<user_id>, sk=DEPOSIT#<deposit_id>.
type BalanceRecord struct {
	Pk           string  `dynamo:"pk"`
	Sk           string  `dynamo:"sk"`
	UserID       string  `dynamo:"user_id"`
	DepositID    string  `dynamo:"deposit_id"`
	Amount       float64 `dynamo:"amount"`
	BalanceAfter float64 `dynamo:"balance_after"`
	CreatedAt    string  `dynamo:"created_at"`
}
