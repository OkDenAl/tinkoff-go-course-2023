package userrepo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework10/internal/entities/user"
	"testing"
)

type UserRepoTest struct {
	Name    string
	Body    any
	Expect  any
	IsError bool
}

type CreateBody struct {
	user *user.User
}

type DeleteBody struct {
	userId int64
}

type GetBody struct {
	userId int64
}

type UpdatePasswordBody struct {
	userId  int64
	newPass string
}

type UpdateNickBody struct {
	userId  int64
	newNick string
}

func TestUserRepo(t *testing.T) {
	tests := []UserRepoTest{
		{
			Name: "create user ok", Body: CreateBody{user: &user.User{Password: "pass1", Email: "email1", Nickname: "nick1"}},
			Expect: int64(1),
		},
		{
			Name: "delete user ok", Body: DeleteBody{userId: int64(0)},
			Expect: nil,
		},
		{
			Name: "update user nickname ok", Body: UpdateNickBody{userId: int64(0), newNick: "new nick"},
			Expect: &user.User{Id: int64(0), Nickname: "new nick"},
		},
		{
			Name: "update user password ok", Body: UpdatePasswordBody{userId: int64(0), newPass: "new pass"},
			Expect: &user.User{Id: int64(0), Password: "new pass"},
		},
		{
			Name: "get user ok", Body: GetBody{userId: int64(0)},
			Expect: &user.User{Id: int64(0), Nickname: "nickname"},
		},
		{
			Name: "get user error", Body: GetBody{userId: int64(-1)},
			Expect: ErrInvalidUserId, IsError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			repo := NewForTest(map[int64]*user.User{
				0: {
					Nickname: "nickname",
					Id:       0,
					Password: "pass",
					Email:    "email",
				},
			}, 1)
			ctx := context.Background()
			switch tc.Body.(type) {
			case CreateBody:
				id, _ := repo.CreateUser(ctx, tc.Body.(CreateBody).user)
				assert.Equal(t, id, tc.Expect)
			case DeleteBody:
				err := repo.DeleteUser(ctx, tc.Body.(DeleteBody).userId)
				assert.Nil(t, err)
			case GetBody:
				u, _ := repo.GetUser(ctx, tc.Body.(GetBody).userId)
				if tc.IsError {
					assert.Error(t, tc.Expect.(error))
				} else {
					assert.Equal(t, u.Nickname, tc.Expect.(*user.User).Nickname)
				}
			case UpdateNickBody:
				u, _ := repo.UpdateNick(ctx, tc.Body.(UpdateNickBody).userId, tc.Body.(UpdateNickBody).newNick)
				assert.Equal(t, u.Nickname, tc.Expect.(*user.User).Nickname)
			case UpdatePasswordBody:
				u, _ := repo.UpdatePassword(ctx, tc.Body.(UpdatePasswordBody).userId, tc.Body.(UpdatePasswordBody).newPass)
				assert.Equal(t, u.Password, tc.Expect.(*user.User).Password)
			}
		})
	}
}
