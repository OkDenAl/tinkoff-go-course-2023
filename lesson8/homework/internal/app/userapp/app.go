package userapp

import (
	"context"
	"homework8/internal/entities/ads"
	"homework8/internal/entities/user"
)

type App interface {
	CreateUser(ctx context.Context, nickname, email, password string) (*user.User, error)
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
	u, err = a.repo.UpdateNick(ctx, id, nickname)
	if err != nil {
		return nil, err
	}
	return u, nil
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
	u, err = a.repo.UpdatePassword(ctx, id, password)
	if err != nil {
		return nil, err
	}
	return u, nil
}
