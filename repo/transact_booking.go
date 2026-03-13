package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/dynamo/v2"

	"booking-be/models"
	"booking-be/view"
)

const (
	TableBookings    = "bookings"
	TableBookedSeats = "booked_seats"
)

func BookSeatsTransaction(
	ctx context.Context,
	db *dynamo.DB,
	booking models.Bookings,
	seats []models.Seat,
) error {
	if len(seats) == 0 {
		return fmt.Errorf("no seats to book")
	}

	bookingTbl := db.Table(TableBookings)
	seatTbl := db.Table(TableBookedSeats)

	bookingRecord := view.BookingDomain2Repo(booking)

	tx := db.WriteTx()
	tx.Put(bookingTbl.Put(bookingRecord))

	now := time.Now().UTC().Format(time.RFC3339)
	available := string(models.SeatStatusAvailable)

	for i := range seats {
		oldUpdatedAt := seats[i].UpdatedAt
		seats[i].BookingID = booking.ID
		seats[i].SeatStatus = models.SeatStatusBooked
		seats[i].UpdatedAt = now
		rec := view.BookedSeatDomain2Repo(seats[i])
		tx.Put(seatTbl.Put(rec).If("'updated_at' = ? AND seat_status = ?", oldUpdatedAt, available))
	}

	if err := tx.Run(ctx); err != nil {
		return fmt.Errorf("book seats transaction: %w", err)
	}
	return nil
}
