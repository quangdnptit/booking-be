package view

import (
	"booking-be/models"
	"booking-be/repomodel"
)

// BookedSeatRepo2Domain maps a repo (persistence) record to the domain model
func BookedSeatRepo2Domain(r repomodel.BookedSeatRecord) models.BookedSeat {
	return models.BookedSeat{
		ID:         r.ID,
		BookingID:  r.BookingID,
		ShowtimeID: r.ShowtimeID,
		SeatID:     r.SeatID,
		Status:     r.Status,
	}
}

// BookedSeatDomain2Repo maps the domain model to a repo (persistence) record
func BookedSeatDomain2Repo(s models.BookedSeat) repomodel.BookedSeatRecord {
	return repomodel.BookedSeatRecord{
		ID:         s.ID,
		BookingID:  s.BookingID,
		ShowtimeID: s.ShowtimeID,
		SeatID:     s.SeatID,
		Status:     s.Status,
	}
}
