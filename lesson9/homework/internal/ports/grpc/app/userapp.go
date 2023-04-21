package app

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"homework9/internal/entities/ads"
	"homework9/internal/entities/user"
	"homework9/internal/ports/grpc/base"
)

type UserService struct {
	repo user.Repository
	base.UnimplementedUserServiceServer
}

func NewUserService(repo user.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) CreateUser(ctx context.Context, req *base.CreateUserRequest) (*base.UserResponse, error) {
	usr := &user.User{
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: req.Password,
	}
	err := user.ValidateUser(usr)
	if err != nil {
		return nil, user.ErrInvalidUserParams
	}
	id, err := us.repo.CreateUser(ctx, usr)
	if err != nil {
		return nil, err
	}
	usr.Id = id
	return &base.UserResponse{
		Id:       usr.Id,
		Nickname: usr.Nickname,
		Email:    usr.Email,
	}, nil
}

func (us *UserService) ChangeNickname(ctx context.Context, req *base.ChangeNicknameRequest) (*base.UserResponse, error) {
	err := user.ValidateUser(&user.User{Nickname: req.Nickname, Password: "mockpass"})
	if err != nil {
		return nil, ads.ErrInvalidAdParams
	}
	usr, err := us.repo.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if usr.Nickname == req.Nickname {
		return &base.UserResponse{
			Id:       usr.Id,
			Nickname: usr.Nickname,
			Email:    usr.Email,
		}, nil
	}
	usr, err = us.repo.UpdateNick(ctx, req.Id, req.Nickname)
	if err != nil {
		return nil, err
	}
	return &base.UserResponse{
		Id:       usr.Id,
		Nickname: usr.Nickname,
		Email:    usr.Email,
	}, nil
}

func (us *UserService) GetUser(ctx context.Context, req *base.GetUserRequest) (*base.UserResponse, error) {
	usr, err := us.repo.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &base.UserResponse{
		Id:       usr.Id,
		Nickname: usr.Nickname,
		Email:    usr.Email,
	}, nil
}

func (us *UserService) DeleteUser(ctx context.Context, req *base.DeleteUserRequest) (*empty.Empty, error) {
	_, err := us.repo.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, us.repo.DeleteUser(ctx, req.Id)
}
