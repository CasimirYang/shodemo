package main

import (
	"log"
	"net"
	"tcpserver/api"

	"google.golang.org/grpc"
	pb "tcpserver/api/protoo"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen:  %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServer(s, &api.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
