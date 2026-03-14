package service

import (
	"context"

	"booking-be/models"
	"booking-be/repo"
)

// ProgramService exposes movie/showtime reads aligned with cms-booking (list + get by id).
type ProgramService struct {
	repo *repo.PostgresProgramRepo
}

func NewProgramService(r *repo.PostgresProgramRepo) *ProgramService {
	return &ProgramService{repo: r}
}

func (s *ProgramService) ListMovies(ctx context.Context) ([]models.MovieResponse, error) {
	return s.repo.ListMovies(ctx)
}

func (s *ProgramService) GetMovieByID(ctx context.Context, id string) (*models.MovieResponse, error) {
	return s.repo.GetMovieByID(ctx, id)
}

func (s *ProgramService) ListShowtimes(ctx context.Context) ([]models.ShowtimeResponse, error) {
	return s.repo.ListShowtimes(ctx)
}

func (s *ProgramService) GetShowtimeByID(ctx context.Context, id string) (*models.ShowtimeResponse, error) {
	return s.repo.GetShowtimeByID(ctx, id)
}

func (s *ProgramService) ListTheaters(ctx context.Context) ([]models.TheaterResponse, error) {
	return s.repo.ListTheaters(ctx)
}

func (s *ProgramService) GetTheaterByID(ctx context.Context, id string) (*models.TheaterResponse, error) {
	return s.repo.GetTheaterByID(ctx, id)
}

func (s *ProgramService) ListRoomsByTheaterID(ctx context.Context, theaterID string) ([]models.RoomResponse, error) {
	return s.repo.ListRoomsByTheaterID(ctx, theaterID)
}
