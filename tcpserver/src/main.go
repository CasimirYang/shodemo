package main

import (
	"context"
	"github.com/CasimirYang/share"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
	"tcpserver/api"
	"tcpserver/api/proto"
)

func main() {
	port := viper.GetString("rpc.port")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		share.SugarLogger.Error(err)
		os.Exit(1)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(share.SugarLogger.Desugar()), logFilter(),
		)),
	)
	proto.RegisterUserServer(s, &api.Server{})
	if err := s.Serve(lis); err != nil {
		share.SugarLogger.Error(err)
		os.Exit(1)
	}
}

//logFilter to log request and response body
func logFilter() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		share.SugarLogger.Infof("reqeust:%s resp: %s", req, resp)
		return resp, err
	}
}
