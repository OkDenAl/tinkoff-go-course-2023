package userrepo

import (
	"context"
	"errors"
	"homework8/internal/entities/user"
	"sync"
)

var (
	ErrInvalidUserId = errors.New("cant find this id in map")
)

type repository struct {
	mu             *sync.Mutex
	userDataById   map[int64]*user.User
	curIdGenerator int64
}

func New() user.Repository {
	var mu sync.Mutex
	return &repository{userDataById: make(map[int64]*user.User), curIdGenerator: 0, mu: &mu}
}

func (r *repository) CreateUser(ctx context.Context, user *user.User) (int64, error) {
	r.mu.Lock()
	user.Id = r.curIdGenerator
	r.userDataById[r.curIdGenerator] = user
	r.curIdGenerator++
	r.mu.Unlock()
	return user.Id, nil
}
func (r *repository) UpdateNick(ctx context.Context, id int64, nick string) (*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userDataById[id].Nickname = nick
	return r.userDataById[id], nil
}
func (r *repository) UpdatePassword(ctx context.Context, id int64, pass string) (*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userDataById[id].Password = pass
	return r.userDataById[id], nil
}

func (r *repository) GetUser(ctx context.Context, id int64) (*user.User, error) {
	r.mu.Lock()
	data, ok := r.userDataById[id]
	r.mu.Unlock()
	if !ok {
		return nil, ErrInvalidUserId
	}
	return data, nil
}
