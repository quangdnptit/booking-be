package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	pbBooking "booking-be/gen/booking"
	pbRoom "booking-be/gen/room"
	
	"booking-be/internal/service/booking"
	"booking-be/internal/service/room"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startGRPCServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port %v: %v", port, err)
	}
	
	s := grpc.NewServer()
	
	// Create service implementation instances
	bookingSvc := booking.NewServer()
	roomSvc := room.NewServer()
	
	// Register services with the gRPC server
	pbBooking.RegisterBookingServiceServer(s, bookingSvc)
	pbRoom.RegisterRoomServiceServer(s, roomSvc)
	
	// Register reflection service on gRPC server to allow using grpcurl structure discovery
	reflection.Register(s)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func startHTTPServer(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to the Go web server powered by native net/http!\n")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK\n")
	})

	log.Printf("HTTP server listening on %v", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("failed to serve HTTP on port %v: %v", port, err)
	}
}

func main() {
	grpcPort := ":50051"
	httpPort := ":8080"

	log.Println("Starting dual HTTP/gRPC application...")

	// Start the gRPC server in a separate background goroutine
	go startGRPCServer(grpcPort)

	// Start the HTTP server. This will block the main goroutine.
	startHTTPServer(httpPort)
}
