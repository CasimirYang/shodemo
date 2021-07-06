package rpc

import (
	"context"
	"tcpserver/domain"
	"tcpserver/rpc/proto"
	"tcpserver/service"

	"github.com/CasimirYang/share"
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
		return &proto.UserInfoReply{Code: share.SystemError, UserInfo: nil}, nil
	} else if err == domain.ErrNoData {
		return &proto.UserInfoReply{Code: share.LoginFailError, UserInfo: nil}, nil
	} else if user != nil {
		userRes := proto.UserInfo{Uid: user.Id, UserName: user.UserName, NickName: user.NickName, Password: user.Password, Profile: user.Profile}
		return &proto.UserInfoReply{Code: share.Success, UserInfo: &userRes}, nil
	}
	return &proto.UserInfoReply{Code: share.Success, UserInfo: nil}, nil
}
