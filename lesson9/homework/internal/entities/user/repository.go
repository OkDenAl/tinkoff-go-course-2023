package user

import "context"

type Repository interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, user *User) (int64, error)
	UpdateNick(ctx context.Context, id int64, nick string) (*User, error)
	UpdatePassword(ctx context.Context, id int64, pass string) (*User, error)
}
