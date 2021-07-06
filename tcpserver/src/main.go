package main

import (
	"github.com/CasimirYang/share"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
	"tcpserver/rpc"
	"tcpserver/rpc/proto"
)

func main() {
	port := viper.GetString("rpc.port")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		share.SugarLogger.Error(err)
		os.Exit(1)
	}
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(share.SugarLogger.Desugar()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(share.SugarLogger.Desugar()),
		)),
	)

	proto.RegisterUserServer(s, &rpc.Server{})
	if err := s.Serve(lis); err != nil {
		share.SugarLogger.Error(err)
		os.Exit(1)
	}
}
