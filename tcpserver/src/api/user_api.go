package api

import (
	"context"
	"git.garena.com/jinghua.yang/entry-task-common/proto"
	"tcpserver/domain"
	"tcpserver/service"

	commonCode "git.garena.com/jinghua.yang/entry-task-common/code"
)

type Server struct {
	proto.UnimplementedUserServer
}

func (s *Server) Login(_ context.Context, in *proto.LoginRequest) (resp *proto.UserInfoReply, err error) {
	user, err := service.Login(in.GetUserName(), in.GetPassword())
	return convertReply(err, user)
}

func (s *Server) GetUser(_ context.Context, in *proto.GetUserRequest) (*proto.UserInfoReply, error) {
	user, err := service.GetUser(in.GetUid())
	return convertReply(err, user)
}

func (s *Server) UpdateUser(_ context.Context, in *proto.UpdateUserRequest) (*proto.UserInfoReply, error) {
	err := service.UpdateUserByUid(in.GetUid(), in.GetNickName(), in.GetProfile())
	return convertReply(err, nil)
}

func convertReply(err error, user *domain.UserDO) (*proto.UserInfoReply, error) {
	if err == domain.ErrSystem {
		return &proto.UserInfoReply{Code: commonCode.SystemError, UserInfo: nil}, nil
	} else if err == domain.ErrNoData {
		return &proto.UserInfoReply{Code: commonCode.LoginFailError, UserInfo: nil}, nil
	} else if user != nil {
		userRes := proto.UserInfo{Uid: user.Id, UserName: user.UserName, NickName: user.NickName, Password: user.Password, Profile: user.Profile}
		return &proto.UserInfoReply{Code: commonCode.Success, UserInfo: &userRes}, nil
	}
	return &proto.UserInfoReply{Code: commonCode.Success, UserInfo: nil}, nil
}
