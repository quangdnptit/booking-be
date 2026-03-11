package booking

import (
	"context"
	pb "booking-be/gen/booking"
	"log"
	
	"github.com/google/uuid"
)

type Server struct {
	pb.UnimplementedBookingServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	log.Printf("CreateBooking called for User: %s, Room: %s", req.GetUserId(), req.GetRoomId())
	
	// Example business logic
	bookingID := uuid.New().String()
	
	return &pb.CreateBookingResponse{
		BookingId: bookingID,
		Status:    "CONFIRMED",
	}, nil
}

func (s *Server) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	log.Printf("GetBooking called for ID: %s", req.GetBookingId())
	
	// Example mocked response
	return &pb.GetBookingResponse{
		BookingId: req.GetBookingId(),
		RoomId:    "room-123",
		UserId:    "user-456",
		StartDate: "2024-01-01",
		EndDate:   "2024-01-05",
		Status:    "CONFIRMED",
	}, nil
}
