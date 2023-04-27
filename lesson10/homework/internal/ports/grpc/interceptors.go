package grpc

import (
	"context"
	"google.golang.org/grpc"
	"homework10/pkg/logger"
	"time"
)

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	log := logger.InitLog()
	log.Info("RequestBody: ", req)

	h, err := handler(ctx, req)

	log.Infof("Request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)
	return h, err
}
