package repo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"booking-be/models"
	"booking-be/repomodel"
	"booking-be/view"
)

const (
	userBookingsIndex     = "user-bookings-index"
	showtimeBookingsIndex = "showtime-bookings-index"
	statusBookingsIndex   = "status-bookings-index"
)

// BookingRepo defines operations for showtime bookings in DynamoDB.
// All methods use the domain model (models.Bookings).
type BookingRepo interface {
	GetByID(ctx context.Context, id string) (*models.Bookings, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Bookings, error)
	GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.Bookings, error)
	GetByStatus(ctx context.Context, status string) ([]models.Bookings, error)
	Create(ctx context.Context, booking models.Bookings) error
	Update(ctx context.Context, booking models.Bookings) error
	UpdateStatus(ctx context.Context, id, status string) error
}

// DynamoBookingRepo implements BookingRepo using DynamoDB
type DynamoBookingRepo struct {
	client *dynamodb.Client
	table  string
}

// NewDynamoBookingRepo creates a new DynamoDB-backed booking repo
func NewDynamoBookingRepo(client *dynamodb.Client, tableName string) *DynamoBookingRepo {
	if tableName == "" {
		tableName = "bookings"
	}
	return &DynamoBookingRepo{client: client, table: tableName}
}

// GetByID returns a booking by id (partition key)
func (r *DynamoBookingRepo) GetByID(ctx context.Context, id string) (*models.Bookings, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get booking: %w", err)
	}
	if out.Item == nil {
		return nil, nil
	}
	var rec repomodel.BookingRecord
	if err := attributevalue.UnmarshalMap(out.Item, &rec); err != nil {
		return nil, fmt.Errorf("unmarshal booking: %w", err)
	}
	domain := view.BookingRepo2Domain(rec)
	return &domain, nil
}

// GetByUserID queries by user_id using user-bookings-index (user_id HASH, created_at RANGE)
func (r *DynamoBookingRepo) GetByUserID(ctx context.Context, userID string) ([]models.Bookings, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String(userBookingsIndex),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query by user_id: %w", err)
	}
	return unmarshalBookingsToDomain(out.Items)
}

// GetByShowtimeID queries by showtime_id using showtime-bookings-index
func (r *DynamoBookingRepo) GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.Bookings, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String(showtimeBookingsIndex),
		KeyConditionExpression: aws.String("showtime_id = :sid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sid": &types.AttributeValueMemberS{Value: showtimeID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query by showtime_id: %w", err)
	}
	return unmarshalBookingsToDomain(out.Items)
}

// GetByStatus queries by status using status-bookings-index
func (r *DynamoBookingRepo) GetByStatus(ctx context.Context, status string) ([]models.Bookings, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String(statusBookingsIndex),
		KeyConditionExpression: aws.String("status = :st"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":st": &types.AttributeValueMemberS{Value: status},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query by status: %w", err)
	}
	return unmarshalBookingsToDomain(out.Items)
}

// Create inserts a new booking
func (r *DynamoBookingRepo) Create(ctx context.Context, booking models.Bookings) error {
	rec := view.BookingDomain2Repo(booking)
	item, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return fmt.Errorf("marshal booking: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put booking: %w", err)
	}
	return nil
}

// Update overwrites the booking item
func (r *DynamoBookingRepo) Update(ctx context.Context, booking models.Bookings) error {
	rec := view.BookingDomain2Repo(booking)
	item, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return fmt.Errorf("marshal booking: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("update booking: %w", err)
	}
	return nil
}

// UpdateStatus updates only the status attribute
func (r *DynamoBookingRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET #status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: status},
		},
	})
	if err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}
	return nil
}

func unmarshalBookingsToDomain(items []map[string]types.AttributeValue) ([]models.Bookings, error) {
	if len(items) == 0 {
		return nil, nil
	}
	var records []repomodel.BookingRecord
	if err := attributevalue.UnmarshalListOfMaps(items, &records); err != nil {
		return nil, err
	}
	out := make([]models.Bookings, len(records))
	for i := range records {
		out[i] = view.BookingRepo2Domain(records[i])
	}
	return out, nil
}
