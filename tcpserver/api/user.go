package api

import (
	"context"
	"tcpserver/domain"
	"tcpserver/service"

	"github.com/CasimirYang/share"
	pb "tcpserver/api/proto"
)

type Server struct {
	pb.UnimplementedUserServer
}

func (s *Server) Login(_ context.Context, in *pb.LoginRequest) (resp *pb.UserInfoReply, err error) {
	user, err := service.Login(in.GetUserName(), in.GetPassword())
	return convertReply(err, user)
}

func (s *Server) GetUser(_ context.Context, in *pb.GetUserRequest) (*pb.UserInfoReply, error) {
	user, err := service.GetUser(in.GetUid())
	return convertReply(err, user)
}

func (s *Server) UpdateUser(_ context.Context, in *pb.UpdateUserRequest) (*pb.UserInfoReply, error) {
	err := service.UpdateUserByUid(in.GetUid(), in.GetNickName(), in.GetProfile())
	return convertReply(err, nil)
}

func convertReply(err error, user *domain.UserDO) (*pb.UserInfoReply, error) {
	if err == domain.ErrSystem {
		return &pb.UserInfoReply{Code: share.SystemError, UserInfo: nil}, nil
	} else if err == domain.ErrNoData {
		return &pb.UserInfoReply{Code: share.NoDataError, UserInfo: nil}, nil
	} else if user != nil {
		userRes := pb.UserInfo{Uid: user.Id, UserName: user.UserName, NickName: user.NickName, Password: user.Password, Profile: user.Profile}
		return &pb.UserInfoReply{Code: share.Success, UserInfo: &userRes}, nil
	}
	return &pb.UserInfoReply{Code: share.Success, UserInfo: nil}, nil
}
