package userrepo

import (
	"context"
	"errors"
	"homework8/internal/entities/user"
	"log"
)

var (
	ErrInvalidUserId = errors.New("cant find this id in map")
)

type repository struct {
	userDataById   map[int64]*user.User
	curIdGenerator int64
}

func New() user.Repository {
	return &repository{userDataById: make(map[int64]*user.User), curIdGenerator: 0}
}

func (r *repository) CreateUser(ctx context.Context, user *user.User) (int64, error) {
	user.Id = r.curIdGenerator
	r.userDataById[r.curIdGenerator] = user
	r.curIdGenerator++
	return user.Id, nil
}
func (r *repository) UpdateNick(ctx context.Context, id int64, nick string) (*user.User, error) {
	r.userDataById[id].Nickname = nick
	return r.userDataById[id], nil
}
func (r *repository) UpdatePassword(ctx context.Context, id int64, pass string) (*user.User, error) {
	r.userDataById[id].Password = pass
	return r.userDataById[id], nil
}

func (r *repository) GetUser(ctx context.Context, id int64) (*user.User, error) {
	log.Println(r.userDataById)
	data, ok := r.userDataById[id]
	if !ok {
		return nil, ErrInvalidUserId
	}
	return data, nil
}
