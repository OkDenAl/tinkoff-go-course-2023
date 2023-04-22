package grpc

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"homework9/internal/entities/ads"
	"homework9/internal/entities/user"
	"homework9/internal/ports/grpc/app"
	"homework9/internal/ports/grpc/base"
)

func NewGrpcServer(adRepo ads.Repository, userRepo user.Repository) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		),
	)
	base.RegisterAdServiceServer(server, app.NewAdService(adRepo, userRepo))
	base.RegisterUserServiceServer(server, app.NewUserService(userRepo))
	return server
}
