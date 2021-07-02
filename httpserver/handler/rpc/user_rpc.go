package rpc

import (
	"context"
	"google.golang.org/grpc"
	pb "httpserver/handler/rpc/proto"
	"log"
	"time"
)

const address = "localhost:50051"

var userClient pb.UserClient

func init() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close() todo
	userClient = pb.NewUserClient(conn)
}

func Login(userName, password string) (*pb.UserInfoReply, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return userClient.Login(ctx, &pb.LoginRequest{UserName: userName, Password: password})
}

func GetUser(uid int64) (*pb.UserInfoReply, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return userClient.GetUser(ctx, &pb.GetUserRequest{Uid: uid})
}

func EditUser(uid int64, nickName, profile *string) (*pb.UserInfoReply, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var updateUserRequest *pb.UpdateUserRequest
	if nickName != nil {
		updateUserRequest = &pb.UpdateUserRequest{Uid: uid, NickName: *nickName}
	} else if profile != nil {
		updateUserRequest = &pb.UpdateUserRequest{Uid: uid, Profile: *profile}
	}
	return userClient.UpdateUser(ctx, updateUserRequest)
}
