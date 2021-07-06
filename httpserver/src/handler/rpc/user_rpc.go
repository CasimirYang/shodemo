package rpc

import (
	"context"
	"github.com/CasimirYang/share"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"httpserver/handler/rpc/proto"
	"time"
)

var userClient proto.UserClient

func init() {
	address := viper.GetString("user.rpcUrl")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		share.SugarLogger.Error(err.Error())
	}
	userClient = proto.NewUserClient(conn)
}

func Login(userName, password string) (*proto.UserInfoReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return userClient.Login(ctx, &proto.LoginRequest{UserName: userName, Password: password})
}

func GetUser(uid int64) (*proto.UserInfoReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return userClient.GetUser(ctx, &proto.GetUserRequest{Uid: uid})
}

func EditUser(uid int64, nickName, profile *string) (*proto.UserInfoReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var updateUserRequest *proto.UpdateUserRequest
	if nickName != nil {
		updateUserRequest = &proto.UpdateUserRequest{Uid: uid, NickName: *nickName}
	} else if profile != nil {
		updateUserRequest = &proto.UpdateUserRequest{Uid: uid, Profile: *profile}
	}
	return userClient.UpdateUser(ctx, updateUserRequest)
}
