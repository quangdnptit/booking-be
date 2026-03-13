package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/dynamo/v2"

	"booking-be/models"
	"booking-be/repomodel"
	"booking-be/view"
)

// SeatRepo defines operations for booked seats in DynamoDB.
// Table key: pk=showtime_id, sk=seat_key. GSI: booking-seats-index.
type SeatRepo interface {
	GetByShowtimeIDAndSeatKeys(ctx context.Context, showtimeID string, seatKeys []string) ([]models.Seat, error)
	GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.Seat, error)
	Update(ctx context.Context, seat models.Seat) error
	UpdateStatusByKey(ctx context.Context, showtimeID, seatKey, status string) error
	GenerateSeats(ctx context.Context, seats []models.Seat) error
}

// DynamoBookedSeatRepo implements SeatRepo via guregu/dynamo (AWS SDK v2).
type DynamoBookedSeatRepo struct {
	table dynamo.Table
}

// NewDynamoBookedSeatRepo creates a repo backed by the given dynamo.DB.
func NewDynamoBookedSeatRepo(db *dynamo.DB) *DynamoBookedSeatRepo {
	return &DynamoBookedSeatRepo{table: db.Table(TableBookedSeats)}
}

func (r *DynamoBookedSeatRepo) GetByShowtimeIDAndSeatKeys(
	ctx context.Context, showtimeID string, seatKeys []string,
) ([]models.Seat, error) {
	if len(seatKeys) == 0 {
		return nil, nil
	}
	batchGetQuery := r.table.Batch("showtime_id", "seat_key").Get(dynamo.Keys{showtimeID, seatKeys[0]})
	for _, sk := range seatKeys[1:] {
		batchGetQuery = batchGetQuery.And(dynamo.Keys{showtimeID, sk})
	}
	var records []repomodel.BookedSeatRecord
	if err := batchGetQuery.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("batch get booked seats: %w", err)
	}
	out := make([]models.Seat, len(records))
	for i := range records {
		out[i] = view.BookedSeatRepo2Domain(records[i])
	}
	return out, nil
}

func (r *DynamoBookedSeatRepo) GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.Seat, error) {
	var records []repomodel.BookedSeatRecord
	err := r.table.Get("showtime_id", showtimeID).All(ctx, &records)
	if err != nil {
		return nil, fmt.Errorf("query by showtime_id: %w", err)
	}
	return bookedRecordsToSeats(records), nil
}

func (r *DynamoBookedSeatRepo) Update(ctx context.Context, seat models.Seat) error {
	rec := view.BookedSeatDomain2Repo(seat)
	newTime := time.Now().UTC().Format(time.RFC3339)
	err := r.table.Update("showtime_id", rec.ShowtimeID).
		Range("seat_key", rec.SeatKey).
		Set("status", rec.Status).
		Set("updated_at", newTime).
		If("'updated_at' = ?", rec.UpdatedAt).
		Run(ctx)
	if err != nil {
		return fmt.Errorf("seat already modified")
	}
	return nil
}

func (r *DynamoBookedSeatRepo) UpdateStatusByKey(ctx context.Context, showtimeID, seatKey, status string) error {
	err := r.table.Update("showtime_id", showtimeID).
		Range("seat_key", seatKey).
		Set("status", status).
		Run(ctx)
	if err != nil {
		return fmt.Errorf("update booked seat status: %w", err)
	}
	return nil
}

func (r *DynamoBookedSeatRepo) GenerateSeats(ctx context.Context, seats []models.Seat) error {
	if len(seats) == 0 {
		return nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	recs := make([]interface{}, 0, len(seats))
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
		recs = append(recs, view.BookedSeatDomain2Repo(s))
	}
	_, err := r.table.Batch("showtime_id", "seat_key").Write().Put(recs...).Run(ctx)
	if err != nil {
		return fmt.Errorf("batch write seats: %w", err)
	}
	return nil
}

func bookedRecordsToSeats(records []repomodel.BookedSeatRecord) []models.Seat {
	if len(records) == 0 {
		return nil
	}
	out := make([]models.Seat, len(records))
	for i := range records {
		out[i] = view.BookedSeatRepo2Domain(records[i])
	}
	return out
}
