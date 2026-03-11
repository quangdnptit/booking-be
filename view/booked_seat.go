package view

import (
	"booking-be/models"
	"booking-be/repomodel"
)

// BookedSeatRepo2Domain maps a repo (persistence) record to the domain model
func BookedSeatRepo2Domain(r repomodel.BookedSeatRecord) models.BookedSeat {
	return models.BookedSeat{
		BookingID:  r.BookingID,
		ShowtimeID: r.ShowtimeID,
		SeatKey:    r.SeatKey,
		Status:     r.Status,
	}
}

// BookedSeatDomain2Repo maps the domain model to a repo (persistence) record
func BookedSeatDomain2Repo(s models.BookedSeat) repomodel.BookedSeatRecord {
	return repomodel.BookedSeatRecord{

		BookingID:  s.BookingID,
		ShowtimeID: s.ShowtimeID,
		SeatKey:    s.SeatKey,
		Status:     s.Status,
	}
}
