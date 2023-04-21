package tests

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
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
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")
	t.Cleanup(func() {
		conn.Close()
	})
	return conn
}

func TestGRRPCCreateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	client := base.NewUserServiceClient(conn)
	res, err := client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")
	assert.Zero(t, res.Id)
	assert.Equal(t, "Oleg", res.Nickname)
	assert.Equal(t, "oleg@tinkoff.com", res.Email)
}

func TestGRRPCGetUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	client := base.NewUserServiceClient(conn)
	_, err := client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")
	_, err = client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg2", Email: "oleg2@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	_, _ = client.ChangeNickname(ctx, &base.ChangeNicknameRequest{Id: int64(0), Nickname: "n"})
	assert.Error(t, user.ErrInvalidUserParams)

	res, err := client.ChangeNickname(ctx, &base.ChangeNicknameRequest{Id: int64(0), Nickname: "Denis"})
	assert.NoError(t, err)

	assert.Zero(t, res.Id)
	assert.Equal(t, "Denis", res.Nickname)
	assert.Equal(t, "oleg@tinkoff.com", res.Email)
}

func TestGRRPCDeleteUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	client := base.NewUserServiceClient(conn)
	_, err := client.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	_, err = client.DeleteUser(ctx, &base.DeleteUserRequest{Id: int64(0)})
	assert.NoError(t, err)
	_, _ = client.DeleteUser(ctx, &base.DeleteUserRequest{Id: int64(0)})
	assert.Error(t, userrepo.ErrInvalidUserId)
}

func TestGRRPCCreateAd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
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

func TestGRRPCGetAdById(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	clientUser := base.NewUserServiceClient(conn)
	clientAd := base.NewAdServiceClient(conn)
	_, err := clientUser.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad", Text: "tester", UserId: 0})
	assert.NoError(t, err)
	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad1", Text: "tester1", UserId: 0})
	assert.NoError(t, err)

	res, err := clientAd.GetAdById(ctx, &base.GetAdByIdRequest{AdId: 0})
	assert.NoError(t, err)
	assert.Zero(t, res.Id)
	assert.Equal(t, res.Title, "tester ad")
	assert.Equal(t, res.Text, "tester")
	assert.Equal(t, res.AuthorId, int64(0))
	assert.Equal(t, res.CreationDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, res.UpdateDate, "")
	assert.False(t, res.Published)
}

func TestGRRPCGetAdByTitle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	clientUser := base.NewUserServiceClient(conn)
	clientAd := base.NewAdServiceClient(conn)
	_, err := clientUser.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad", Text: "tester", UserId: 0})
	assert.NoError(t, err)
	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad1", Text: "tester1", UserId: 0})
	assert.NoError(t, err)

	res, err := clientAd.GetAdByTitle(ctx, &base.GetAdByTitleRequest{Title: "tester"})
	assert.NoError(t, err)
	assert.Len(t, res.List, 2)
}

func TestGRRPCChangeAdStatus(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	clientUser := base.NewUserServiceClient(conn)
	clientAd := base.NewAdServiceClient(conn)
	_, err := clientUser.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	res, err := clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad", Text: "tester", UserId: 0})
	assert.NoError(t, err)
	assert.Equal(t, res.UpdateDate, "")
	assert.Equal(t, res.Published, false)
	res, err = clientAd.ChangeAdStatus(ctx, &base.ChangeAdStatusRequest{AdId: 0, UserId: 0, Published: true})
	assert.NoError(t, err)
	assert.Equal(t, res.UpdateDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, res.Published, true)
}

func TestGRRPCUpdateAd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	clientUser := base.NewUserServiceClient(conn)
	clientAd := base.NewAdServiceClient(conn)
	_, err := clientUser.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	res, err := clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad", Text: "tester", UserId: 0})
	assert.NoError(t, err)
	assert.Equal(t, res.UpdateDate, "")
	assert.Equal(t, res.Title, "tester ad")
	assert.Equal(t, res.Text, "tester")
	res, err = clientAd.UpdateAd(ctx, &base.UpdateAdRequest{UserId: 0, AdId: 0, Title: "new title", Text: "new text"})
	assert.NoError(t, err)
	assert.Equal(t, res.UpdateDate, time.Now().UTC().Format(time.DateOnly))
	assert.Equal(t, res.Title, "new title")
	assert.Equal(t, res.Text, "new text")
}

func TestGRRPCGetAllWithFilters(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	clientUser := base.NewUserServiceClient(conn)
	clientAd := base.NewAdServiceClient(conn)
	_, err := clientUser.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad", Text: "tester", UserId: 0})
	assert.NoError(t, err)
	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester1 ad", Text: "tester1", UserId: 0})
	assert.NoError(t, err)
	_, err = clientAd.ChangeAdStatus(ctx, &base.ChangeAdStatusRequest{AdId: 1, UserId: 0, Published: true})
	assert.NoError(t, err)
	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester1 ad", Text: "tester1", UserId: 0})
	assert.NoError(t, err)

	res, err := clientAd.ListAds(ctx, &base.Filters{AuthorId: "1"})
	assert.NoError(t, err)
	assert.Len(t, res.List, 0)
	res, err = clientAd.ListAds(ctx, &base.Filters{AuthorId: "0", Status: "published"})
	assert.NoError(t, err)
	assert.Len(t, res.List, 1)
	assert.Equal(t, res.List[0].Title, "tester1 ad")
	res, err = clientAd.ListAds(ctx, &base.Filters{AuthorId: "0", Status: "unpublished", Date: time.Now().UTC().Format(time.DateOnly)})
	assert.NoError(t, err)
	assert.Len(t, res.List, 2)
}

func TestGRRPCDeleteAd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})
	conn := server(ctx, t)
	clientUser := base.NewUserServiceClient(conn)
	clientAd := base.NewAdServiceClient(conn)
	_, err := clientUser.CreateUser(ctx, &base.CreateUserRequest{Nickname: "Oleg", Email: "oleg@tinkoff.com", Password: "easyhw"})
	assert.NoError(t, err, "client.GetUser")

	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad", Text: "tester", UserId: 0})
	assert.NoError(t, err)
	_, err = clientAd.CreateAd(ctx, &base.CreateAdRequest{Title: "tester ad1", Text: "tester1", UserId: 0})
	assert.NoError(t, err)

	_, err = clientAd.DeleteAd(ctx, &base.DeleteAdRequest{AdId: 1, AuthorId: 0})
	assert.NoError(t, err)
	_, _ = clientAd.DeleteAd(ctx, &base.DeleteAdRequest{AdId: 1, AuthorId: 0})
	assert.Error(t, adrepo.ErrInvalidAdId)
}
