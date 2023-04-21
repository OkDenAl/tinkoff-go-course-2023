package userapp

import (
	"context"
	"homework9/internal/entities/ads"
	"homework9/internal/entities/user"
)

type App interface {
	CreateUser(ctx context.Context, nickname, email, password string) (*user.User, error)
	GetUser(ctx context.Context, id int64) (*user.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ChangeNickname(ctx context.Context, id int64, nickname string) (*user.User, error)
	UpdatePassword(ctx context.Context, id int64, password string) (*user.User, error)
}

type app struct {
	repo user.Repository
}

func NewApp(repo user.Repository) App {
	return app{repo: repo}
}

func (a app) CreateUser(ctx context.Context, nickname, email, password string) (*user.User, error) {
	u := &user.User{
		Nickname: nickname,
		Email:    email,
		Password: password,
	}
	err := user.ValidateUser(u)
	if err != nil {
		return nil, user.ErrInvalidUserParams
	}
	id, err := a.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	u.Id = id
	return u, nil
}

func (a app) ChangeNickname(ctx context.Context, id int64, nickname string) (*user.User, error) {
	err := user.ValidateUser(&user.User{Nickname: nickname, Password: "mockpass"})
	if err != nil {
		return nil, ads.ErrInvalidAdParams
	}
	u, err := a.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if u.Nickname == nickname {
		return u, nil
	}
	return a.repo.UpdateNick(ctx, id, nickname)
}

func (a app) UpdatePassword(ctx context.Context, id int64, password string) (*user.User, error) {
	err := user.ValidateUser(&user.User{Nickname: "mocknick", Password: password})
	if err != nil {
		return nil, ads.ErrInvalidAdParams
	}
	u, err := a.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if u.Password == password {
		return u, nil
	}
	return a.repo.UpdatePassword(ctx, id, password)
}

func (a app) DeleteUser(ctx context.Context, id int64) error {
	_, err := a.repo.GetUser(ctx, id)
	if err != nil {
		return err
	}
	return a.repo.DeleteUser(ctx, id)
}

func (a app) GetUser(ctx context.Context, id int64) (*user.User, error) {
	return a.repo.GetUser(ctx, id)
}
