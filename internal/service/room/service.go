package room

import (
	"context"
	pb "booking-be/gen/room"
	"log"
)

type Server struct {
	pb.UnimplementedRoomServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ListRooms(ctx context.Context, req *pb.ListRoomsRequest) (*pb.ListRoomsResponse, error) {
	log.Println("ListRooms called")
	
	rooms := []*pb.Room{
		{
			Id:            "room-1",
			Name:          "Deluxe Suite",
			Capacity:      2,
			PricePerNight: 150.00,
		},
		{
			Id:            "room-2",
			Name:          "Family Room",
			Capacity:      4,
			PricePerNight: 200.00,
		},
	}
	
	return &pb.ListRoomsResponse{
		Rooms: rooms,
	}, nil
}

func (s *Server) GetRoomAvailability(ctx context.Context, req *pb.GetRoomAvailabilityRequest) (*pb.GetRoomAvailabilityResponse, error) {
	log.Printf("GetRoomAvailability called for Room: %s", req.GetRoomId())
	
	// Fake logic: Always assume available for this mock
	return &pb.GetRoomAvailabilityResponse{
		IsAvailable: true,
	}, nil
}
