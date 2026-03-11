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
	bookingSeatsIndex  = "booking-seats-index"
	showtimeSeatsIndex = "showtime-seats-index"
	statusSeatsIndex   = "status-seats-index"
)

// BookedSeatRepo defines operations for booked seats in DynamoDB.
// All methods use the domain model (models.BookedSeat).
type BookedSeatRepo interface {
	GetByID(ctx context.Context, id string) (*models.BookedSeat, error)
	GetByBookingID(ctx context.Context, bookingID string) ([]models.BookedSeat, error)
	GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.BookedSeat, error)
	GetByStatus(ctx context.Context, status string) ([]models.BookedSeat, error)
	Create(ctx context.Context, seat models.BookedSeat) error
	Update(ctx context.Context, seat models.BookedSeat) error
	UpdateStatus(ctx context.Context, id, status string) error
}

// DynamoBookedSeatRepo implements BookedSeatRepo using DynamoDB
type DynamoBookedSeatRepo struct {
	client *dynamodb.Client
	table  string
}

// NewDynamoBookedSeatRepo creates a new DynamoDB-backed booked seat repo
func NewDynamoBookedSeatRepo(client *dynamodb.Client, tableName string) *DynamoBookedSeatRepo {
	if tableName == "" {
		tableName = "booked_seats"
	}
	return &DynamoBookedSeatRepo{client: client, table: tableName}
}

// GetByID returns a booked seat by id (partition key)
func (r *DynamoBookedSeatRepo) GetByID(ctx context.Context, id string) (*models.BookedSeat, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get booked seat: %w", err)
	}
	if out.Item == nil {
		return nil, nil
	}
	var rec repomodel.BookedSeatRecord
	if err := attributevalue.UnmarshalMap(out.Item, &rec); err != nil {
		return nil, fmt.Errorf("unmarshal booked seat: %w", err)
	}
	domain := view.BookedSeatRepo2Domain(rec)
	return &domain, nil
}

// GetByBookingID queries by booking_id using booking-seats-index (booking_id HASH, seat_id RANGE)
func (r *DynamoBookedSeatRepo) GetByBookingID(ctx context.Context, bookingID string) ([]models.BookedSeat, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String(bookingSeatsIndex),
		KeyConditionExpression: aws.String("booking_id = :bid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":bid": &types.AttributeValueMemberS{Value: bookingID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query by booking_id: %w", err)
	}
	return unmarshalBookedSeatsToDomain(out.Items)
}

// GetByShowtimeID queries by showtime_id using showtime-seats-index
func (r *DynamoBookedSeatRepo) GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.BookedSeat, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String(showtimeSeatsIndex),
		KeyConditionExpression: aws.String("showtime_id = :sid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sid": &types.AttributeValueMemberS{Value: showtimeID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query by showtime_id: %w", err)
	}
	return unmarshalBookedSeatsToDomain(out.Items)
}

// GetByStatus queries by status using status-seats-index
func (r *DynamoBookedSeatRepo) GetByStatus(ctx context.Context, status string) ([]models.BookedSeat, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String(statusSeatsIndex),
		KeyConditionExpression: aws.String("status = :st"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":st": &types.AttributeValueMemberS{Value: status},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query by status: %w", err)
	}
	return unmarshalBookedSeatsToDomain(out.Items)
}

// Create inserts a new booked seat
func (r *DynamoBookedSeatRepo) Create(ctx context.Context, seat models.BookedSeat) error {
	rec := view.BookedSeatDomain2Repo(seat)
	item, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return fmt.Errorf("marshal booked seat: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put booked seat: %w", err)
	}
	return nil
}

// Update overwrites the booked seat item
func (r *DynamoBookedSeatRepo) Update(ctx context.Context, seat models.BookedSeat) error {
	rec := view.BookedSeatDomain2Repo(seat)
	item, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return fmt.Errorf("marshal booked seat: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("update booked seat: %w", err)
	}
	return nil
}

// UpdateStatus updates only the status attribute
func (r *DynamoBookedSeatRepo) UpdateStatus(ctx context.Context, id, status string) error {
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
		return fmt.Errorf("update booked seat status: %w", err)
	}
	return nil
}

func unmarshalBookedSeatsToDomain(items []map[string]types.AttributeValue) ([]models.BookedSeat, error) {
	if len(items) == 0 {
		return nil, nil
	}
	var records []repomodel.BookedSeatRecord
	if err := attributevalue.UnmarshalListOfMaps(items, &records); err != nil {
		return nil, err
	}
	out := make([]models.BookedSeat, len(records))
	for i := range records {
		out[i] = view.BookedSeatRepo2Domain(records[i])
	}
	return out, nil
}
