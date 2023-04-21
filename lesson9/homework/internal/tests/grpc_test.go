package tests

import (
	"context"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/entities/user"
	"homework9/internal/ports/grpc/app"
	"homework9/internal/ports/grpc/base"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func server(ctx context.Context, t *testing.T) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})
	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})
	userRepo := userrepo.New()
	base.RegisterUserServiceServer(srv, app.NewUserService(userRepo))
	base.RegisterAdServiceServer(srv, app.NewAdService(adrepo.New(), userRepo))
	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	assert.NoError(t, err, "grpc.DialContext")
	t.Cleanup(func() {
		conn.Close()
	})
	return conn
}

func TestGRRPCCreateUser(t *testing.T) {
	ctx := context.Background()
	conn := server(ctx, t)
	client := base.NewUserServiceClient(conn)
	res, err := client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")
	assert.Zero(t, res.Id)
	assert.Equal(t, "Oleg", res.Nickname)
	assert.Equal(t, "oleg@tinkoff.com", res.Email)
}

func TestGRRPCGetUser(t *testing.T) {
	ctx := context.Background()
	conn := server(ctx, t)
	client := base.NewUserServiceClient(conn)
	_, err := client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")
	_, err = client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg2", Email: "oleg2@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	res, err := client.GetUser(ctx, &base.GetUserRequest{Id: int64(0)})
	assert.NoError(t, err)

	assert.Zero(t, res.Id)
	assert.Equal(t, "Oleg", res.Nickname)
	assert.Equal(t, "oleg@tinkoff.com", res.Email)
}

func TestGRRPCChangeNickname(t *testing.T) {
	ctx := context.Background()
	conn := server(ctx, t)
	client := base.NewUserServiceClient(conn)
	_, err := client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")
	_, err = client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg2", Email: "oleg2@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	_, err = client.ChangeNickname(ctx, &base.ChangeNicknameRequest{Id: int64(0), Nickname: "n"})
	assert.Error(t, user.ErrInvalidUserParams)

	res, err := client.ChangeNickname(ctx, &base.ChangeNicknameRequest{Id: int64(0), Nickname: "Denis"})
	assert.NoError(t, err)

	assert.Zero(t, res.Id)
	assert.Equal(t, "Denis", res.Nickname)
	assert.Equal(t, "oleg@tinkoff.com", res.Email)
}

func TestGRRPCDeleteUser(t *testing.T) {
	ctx := context.Background()
	conn := server(ctx, t)
	client := base.NewUserServiceClient(conn)
	_, err := client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	_, err = client.DeleteUser(ctx, &base.DeleteUserRequest{Id: int64(0)})
	assert.NoError(t, err)
	_, err = client.DeleteUser(ctx, &base.DeleteUserRequest{Id: int64(0)})
	assert.Error(t, userrepo.ErrInvalidUserId)
}

func TestGRRPCCreateAd(t *testing.T) {
	ctx := context.Background()
	conn := server(ctx, t)
	clientUser := base.NewUserServiceClient(conn)
	clientAd := base.NewAdServiceClient(conn)
	_, err := clientUser.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	res, err := clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad", Text: "tester", UserId: 0})
	assert.NoError(t, err)

	assert.Zero(t, res.Id)
	assert.Equal(t, res.Title, "tester ad")
	assert.Equal(t, res.Text, "tester")
	assert.Equal(t, res.AuthorId, int64(0))
	assert.Equal(t, res.CreationDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, res.UpdateDate, "")
	assert.False(t, res.Published)
}
