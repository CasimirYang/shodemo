package main

import (
	"fmt"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"log"
	"net"
	"tcpserver/rpc"
	"tcpserver/trace"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	pb "tcpserver/rpc/proto"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		_ = trace.Logger.Error(fmt.Sprintf("failed to listen:  %v", err))
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(ZapInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(ZapInterceptor()),
		)),
	)

	pb.RegisterUserServer(s, &rpc.Server{})
	if err := s.Serve(lis); err != nil {
		_ = trace.Logger.Error(fmt.Sprintf("failed to serve: %v", err))
	}
}

func ZapInterceptor() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)
	return logger
}
