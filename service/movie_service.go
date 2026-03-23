package service

import (
	"context"

	"github.com/quangdnptit/booking-be/models"
	"github.com/quangdnptit/booking-be/repo"
)

// MovieService exposes movie/showtime reads aligned with cms-booking (list + get by id).
type MovieService struct {
	repo *repo.PostgresProgramRepo
}

func NewMovieService(r *repo.PostgresProgramRepo) *MovieService {
	return &MovieService{repo: r}
}

func (s *MovieService) ListMovies(ctx context.Context) ([]models.MovieResponse, error) {
	return s.repo.ListMovies(ctx)
}

func (s *MovieService) GetMovieByID(ctx context.Context, id string) (*models.MovieResponse, error) {
	return s.repo.GetMovieByID(ctx, id)
}

func (s *MovieService) ListShowtimes(ctx context.Context) ([]models.ShowtimeResponse, error) {
	return s.repo.ListShowtimes(ctx)
}

func (s *MovieService) GetShowtimeByID(ctx context.Context, id string) (*models.ShowtimeResponse, error) {
	return s.repo.GetShowtimeByID(ctx, id)
}

func (s *MovieService) ListTheaters(ctx context.Context) ([]models.TheaterResponse, error) {
	return s.repo.ListTheaters(ctx)
}

func (s *MovieService) GetTheaterByID(ctx context.Context, id string) (*models.TheaterResponse, error) {
	return s.repo.GetTheaterByID(ctx, id)
}

func (s *MovieService) ListRoomsByTheaterID(ctx context.Context, theaterID string) ([]models.RoomResponse, error) {
	return s.repo.ListRoomsByTheaterID(ctx, theaterID)
}
