package rpc

import (
	"context"
	commonLog "git.garena.com/jinghua.yang/entry-task-common/log"
	"git.garena.com/jinghua.yang/entry-task-common/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"time"
)

var userClient proto.UserClient

func init() {
	address := viper.GetString("user.rpcUrl")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		commonLog.SugarLogger.Error(err.Error())
	}
	userClient = proto.NewUserClient(conn)
}

func Login(userName, password string) (*proto.UserInfoReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return userClient.Login(ctx, &proto.LoginRequest{UserName: userName, Password: password})
}

func GetUser(uid int64) (*proto.UserInfoReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return userClient.GetUser(ctx, &proto.GetUserRequest{Uid: uid})
}

func EditUser(uid int64, nickName, profile *string) (*proto.UserInfoReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var updateUserRequest *proto.UpdateUserRequest
	if nickName != nil {
		updateUserRequest = &proto.UpdateUserRequest{Uid: uid, NickName: *nickName}
	} else if profile != nil {
		updateUserRequest = &proto.UpdateUserRequest{Uid: uid, Profile: *profile}
	}
	return userClient.UpdateUser(ctx, updateUserRequest)
}
