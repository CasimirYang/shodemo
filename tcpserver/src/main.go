package main

import (
	"fmt"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"tcpserver/rpc"
	"tcpserver/trace"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	pb "tcpserver/rpc/proto"

	"gopkg.in/natefinch/lumberjack.v2"
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

// ZapInterceptor 返回zap.logger实例(把日志写到文件中)
func ZapInterceptor() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  "log/rpc.log",
		MaxSize:   1024, //MB
		LocalTime: true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.NewAtomicLevel(),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	grpc_zap.ReplaceGrpcLoggerV2(logger)
	return logger
}
