package userapp

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app/userapp/mocks"
	"homework10/internal/entities/user"
	"testing"
)

func TestUserService_CreateUserOK(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("CreateUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*user.User")).
		Return(int64(0), nil)

	service := NewApp(repo)

	u, err := service.CreateUser(context.Background(), "test", "test@gmail.com", "password")
	assert.Zero(t, u.Id)
	assert.Equal(t, u.Nickname, "test")
	assert.Equal(t, u.Email, "test@gmail.com")
	assert.Nil(t, err)
}

func TestUserService_CreateUserInvalidParams(t *testing.T) {
	repo := &mocks.Repository{}
	service := NewApp(repo)

	_, _ = service.CreateUser(context.Background(), "", "", "")
	assert.Error(t, user.ErrInvalidUserParams)
}

func TestUserService_GetUserOK(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&user.User{Id: 0}, nil)

	service := NewApp(repo)

	u, err := service.GetUser(context.Background(), int64(0))
	assert.Nil(t, err)
	assert.Equal(t, u.Id, int64(0))
}

func TestUserService_DeleteUserOK(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&user.User{Id: 0}, nil)
	repo.On("DeleteUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(nil)

	service := NewApp(repo)

	err := service.DeleteUser(context.Background(), int64(0))
	assert.Nil(t, err)
}

func TestUserService_DeleteUserInvalidUserId(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&user.User{Id: 0}, userrepo.ErrInvalidUserId)

	service := NewApp(repo)

	_ = service.DeleteUser(context.Background(), int64(0))
	assert.Error(t, userrepo.ErrInvalidUserId)
}

func TestUserService_UpdatePasswordOK(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&user.User{Id: 0}, nil)
	repo.On("UpdatePassword", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(&user.User{Id: 0, Password: "new password"}, nil)

	service := NewApp(repo)

	u, err := service.UpdatePassword(context.Background(), int64(0), "new password")
	assert.Nil(t, err)
	assert.Equal(t, u.Password, "new password")
}

func TestUserService_UpdatePasswordInvalidUserId(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&user.User{Id: 0}, userrepo.ErrInvalidUserId)

	service := NewApp(repo)

	_, _ = service.UpdatePassword(context.Background(), int64(0), "new password")
	assert.Error(t, userrepo.ErrInvalidUserId)
}

func TestUserService_ChangeNicknameOK(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&user.User{Id: 0}, nil)
	repo.On("UpdateNick", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).
		Return(&user.User{Id: 0, Nickname: "new nickname"}, nil)

	service := NewApp(repo)

	u, err := service.ChangeNickname(context.Background(), int64(0), "new nickname")
	assert.Nil(t, err)
	assert.Equal(t, u.Nickname, "new nickname")
}

func TestUserService_ChangeNicknameInvalidUserId(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetUser", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("int64")).
		Return(&user.User{Id: 0}, userrepo.ErrInvalidUserId)

	service := NewApp(repo)

	_, _ = service.ChangeNickname(context.Background(), int64(0), "new nickname")
	assert.Error(t, userrepo.ErrInvalidUserId)
}
