package repo

import (
	"context"
	"fmt"

	dynamo "github.com/guregu/dynamo/v2"

	"booking-be/models"
	"booking-be/repomodel"
	"booking-be/view"
)

const (
	userBookingsIndex     = "user-bookings-index"
	showtimeBookingsIndex = "showtime-bookings-index"
)

type BookingRepo interface {
	GetByID(ctx context.Context, id string) (*models.Bookings, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Bookings, error)
	GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.Bookings, error)
	Create(ctx context.Context, booking models.Bookings) error
	Update(ctx context.Context, booking models.Bookings) error
	UpdateStatus(ctx context.Context, id, status string) error
}

type DynamoBookingRepo struct {
	table dynamo.Table
}

func NewDynamoBookingRepo(db *dynamo.DB) *DynamoBookingRepo {
	return &DynamoBookingRepo{table: db.Table(TableBookings)}
}

func (r *DynamoBookingRepo) GetByID(ctx context.Context, id string) (*models.Bookings, error) {
	var recs []repomodel.BookingRecord
	err := r.table.Get("id", id).Limit(1).All(ctx, &recs)
	if err != nil {
		return nil, fmt.Errorf("get booking: %w", err)
	}
	if len(recs) == 0 {
		return nil, nil
	}
	d := view.BookingRepo2Domain(recs[0])
	return &d, nil
}

func (r *DynamoBookingRepo) GetByUserID(ctx context.Context, userID string) ([]models.Bookings, error) {
	var recs []repomodel.BookingRecord
	err := r.table.Get("user_id", userID).Index(userBookingsIndex).All(ctx, &recs)
	if err != nil {
		return nil, fmt.Errorf("query by user_id: %w", err)
	}
	return bookingsRecordsToDomain(recs), nil
}

func (r *DynamoBookingRepo) GetByShowtimeID(ctx context.Context, showtimeID string) ([]models.Bookings, error) {
	var recs []repomodel.BookingRecord
	err := r.table.Get("showtime_id", showtimeID).Index(showtimeBookingsIndex).All(ctx, &recs)
	if err != nil {
		return nil, fmt.Errorf("query by showtime_id: %w", err)
	}
	return bookingsRecordsToDomain(recs), nil
}

func (r *DynamoBookingRepo) Create(ctx context.Context, booking models.Bookings) error {
	rec := view.BookingDomain2Repo(booking)
	return r.table.Put(rec).Run(ctx)
}

func (r *DynamoBookingRepo) Update(ctx context.Context, booking models.Bookings) error {
	rec := view.BookingDomain2Repo(booking)
	return r.table.Put(rec).Run(ctx)
}

func (r *DynamoBookingRepo) UpdateStatus(ctx context.Context, id, status string) error {
	b, err := r.GetByID(ctx, id)
	if err != nil || b == nil {
		return err
	}
	return r.table.Update("id", b.ID).
		Range("created_at", b.CreatedAt).
		Set("status", status).
		Run(ctx)
}

func bookingsRecordsToDomain(recs []repomodel.BookingRecord) []models.Bookings {
	if len(recs) == 0 {
		return nil
	}
	out := make([]models.Bookings, len(recs))
	for i := range recs {
		out[i] = view.BookingRepo2Domain(recs[i])
	}
	return out
}
