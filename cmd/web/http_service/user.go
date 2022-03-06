package http_service

import (
	"simple_grpc/internal/facade"
	"simple_grpc/proto"
)

var userClient proto.UserServiceClient

func (svc *Service) getUserClient() proto.UserServiceClient {
	if userClient == nil {
		userClient = proto.NewUserServiceClient(svc.client)
	}
	return userClient
}

func (svc *Service) Register(param *facade.RegisterUserRequest) (*facade.RegisterUserResponse, error) {
	_, err := svc.getUserClient().Register(svc.ctx, &proto.RegisterRequest{
		Username:   param.Username,
		Password:   param.Password,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}
	return &facade.RegisterUserResponse{}, nil
}

func (svc *Service) Login(param *facade.LoginUserRequest) (*facade.LoginUserResponse, error) {
	resp, err := svc.getUserClient().Login(svc.ctx, &proto.LoginRequest{
		Username: param.Username,
		Password: param.Password,
	})
	if err != nil {
		return nil, err
	}
	return &facade.LoginUserResponse{SessionID: resp.SessionId}, nil
}

func (svc *Service) GetUser(param *facade.GetUserRequest) (*facade.GetUserResponse, error) {
	userReply, err := svc.getUserClient().GetUser(svc.ctx, &proto.GetUserRequest{
		SessionId: param.SessionId,
	})
	if err != nil {
		return nil, err
	}
	resp := &facade.GetUserResponse{
		Username:   userReply.Username,
		Nickname:   userReply.Nickname,
		ProfilePic: userReply.ProfilePic,
	}
	return resp, nil
}

func (svc *Service) EditUser(param *facade.EditUserRequest) (*facade.EditUserResponse, error) {
	_, err := svc.getUserClient().EditUser(svc.ctx, &proto.EditUserRequest{
		SessionId:  param.SessionId,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}
	resp := &facade.EditUserResponse{}
	return resp, nil
}
