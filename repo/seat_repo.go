package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"booking-be/models"
	"booking-be/repomodel"
	"booking-be/view"
)

const (
	bookingSeatsIndex = "booking-seats-index" // booking_id HASH, seat_key RANGE
)

// batchWriteMaxItems is DynamoDB BatchWriteItem limit per request
const batchWriteMaxItems = 25

// SeatRepo defines operations for booked seats in DynamoDB.
// Table key: pk=showtime_id, sk=seat_key. GSI: booking-seats-index.
type SeatRepo interface {
	GetByShowtimeIDAndSeatKey(ctx context.Context, showtimeID, seatKey string) (*models.Seat, error)
	GetByBookingID(ctx context.Context, bookingID string) ([]models.Seat, error)
	GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.Seat, error)
	Update(ctx context.Context, seat models.Seat) error
	UpdateStatusByKey(ctx context.Context, showtimeID, seatKey, status string) error
	// GenerateSeats converts domain seats to repo records and batch-writes them (Put)
	GenerateSeats(ctx context.Context, seats []models.Seat) error
}

// DynamoBookedSeatRepo implements SeatRepo using DynamoDB
type DynamoBookedSeatRepo struct {
	client *dynamodb.Client
	table  string
}

// NewDynamoBookedSeatRepo creates a new DynamoDB-backed booked seat repo
func NewDynamoBookedSeatRepo(client *dynamodb.Client, tableName string) *DynamoBookedSeatRepo {
	return &DynamoBookedSeatRepo{client: client, table: tableName}
}

// GetByShowtimeIDAndSeatKey returns one item by table pk=showtime_id, sk=seat_key
func (r *DynamoBookedSeatRepo) GetByShowtimeIDAndSeatKey(ctx context.Context, showtimeID, seatKey string) (*models.Seat, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"showtime_id": &types.AttributeValueMemberS{Value: showtimeID},
			"seat_key":    &types.AttributeValueMemberS{Value: seatKey},
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

// GetByBookingID queries GSI booking-seats-index (booking_id HASH, seat_key RANGE)
func (r *DynamoBookedSeatRepo) GetByBookingID(ctx context.Context, bookingID string) ([]models.Seat, error) {
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

// GetByShowtimeID queries by table pk=showtime_id
func (r *DynamoBookedSeatRepo) GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.Seat, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
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

func (r *DynamoBookedSeatRepo) Update(ctx context.Context, seat models.Seat) error {
	rec := view.BookedSeatDomain2Repo(seat)

	key, err := attributevalue.MarshalMap(map[string]string{
		"showtime_id": rec.ShowtimeID,
		"seat_key":    rec.SeatKey,
	})

	if err != nil {
		return err
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:           aws.String(r.table),
		Key:                 key,
		UpdateExpression:    aws.String("SET #status = :status, #updatedAt = :newTime"),
		ConditionExpression: aws.String("#updatedAt = :oldTime"),
		ExpressionAttributeNames: map[string]string{
			"#status":    "status",
			"#updatedAt": "updatedAt",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":  &types.AttributeValueMemberS{Value: rec.Status},
			":newTime": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
			":oldTime": &types.AttributeValueMemberS{Value: rec.UpdatedAt},
		},
	})

	if err != nil {
		return fmt.Errorf("seat already modified")
	}

	return nil
}

// UpdateStatusByKey updates status by table key (pk=showtime_id, sk=seat_key)
func (r *DynamoBookedSeatRepo) UpdateStatusByKey(ctx context.Context, showtimeID, seatKey, status string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"showtime_id": &types.AttributeValueMemberS{Value: showtimeID},
			"seat_key":    &types.AttributeValueMemberS{Value: seatKey},
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

// GenerateSeats maps domain seats to repo records and batch-saves with BatchWriteItem
func (r *DynamoBookedSeatRepo) GenerateSeats(ctx context.Context, seats []models.Seat) error {
	if len(seats) == 0 {
		return nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	writes := make([]types.WriteRequest, 0, len(seats))
	for i := range seats {
		s := seats[i]
		if s.ShowtimeID == "" || s.SeatKey == "" {
			return fmt.Errorf("seat %d: showtime_id and seat_key are required", i)
		}
		if s.CreatedAt == "" {
			s.CreatedAt = now
		}
		if s.UpdatedAt == "" {
			s.UpdatedAt = now
		}

		rec := view.BookedSeatDomain2Repo(s)
		item, err := attributevalue.MarshalMap(rec)

		if err != nil {
			return fmt.Errorf("marshal seat %s/%s: %w", s.ShowtimeID, s.SeatKey, err)
		}
		writes = append(writes, types.WriteRequest{
			PutRequest: &types.PutRequest{Item: item},
		})
	}
	for start := 0; start < len(writes); start += batchWriteMaxItems {
		end := start + batchWriteMaxItems
		if end > len(writes) {
			end = len(writes)
		}
		chunk := writes[start:end]
		if err := r.batchWriteWithRetry(ctx, chunk); err != nil {
			return err
		}
	}
	return nil
}

func (r *DynamoBookedSeatRepo) batchWriteWithRetry(ctx context.Context, writes []types.WriteRequest) error {
	unprocessed := writes
	backoff := 50 * time.Millisecond
	for attempt := 0; attempt < 10 && len(unprocessed) > 0; attempt++ {
		out, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.table: unprocessed,
			},
		})
		if err != nil {
			return fmt.Errorf("batch write seats: %w", err)
		}
		next := out.UnprocessedItems[r.table]
		if len(next) == 0 {
			return nil
		}
		unprocessed = next
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		if backoff < 2*time.Second {
			backoff *= 2
		}
	}
	if len(unprocessed) > 0 {
		return fmt.Errorf("batch write: %d unprocessed items after retries", len(unprocessed))
	}
	return nil
}

func unmarshalBookedSeatsToDomain(items []map[string]types.AttributeValue) ([]models.Seat, error) {
	if len(items) == 0 {
		return nil, nil
	}
	var records []repomodel.BookedSeatRecord
	if err := attributevalue.UnmarshalListOfMaps(items, &records); err != nil {
		return nil, err
	}
	out := make([]models.Seat, len(records))
	for i := range records {
		out[i] = view.BookedSeatRepo2Domain(records[i])
	}
	return out, nil
}
