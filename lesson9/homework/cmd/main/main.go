package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app/adsapp"
	"homework9/internal/app/userapp"
	grpcInterface "homework9/internal/ports/grpc"
	"homework9/internal/ports/httpgin"
	"homework9/pkg/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	httpPort = ":18081"
	grpcPort = ":50055"
)

func main() {
	userRepo := userrepo.New()
	adsRepo := adrepo.New()
	log := logger.InitLog()

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Info("err with listening grpc port")
		panic(err)
	}

	httpServer := httpgin.NewHTTPServer(httpPort, adsapp.NewApp(adsRepo, userRepo), userapp.NewApp(userRepo), log)
	grpcServer := grpcInterface.NewGrpcServer(adsRepo, userRepo)

	g, ctx := errgroup.WithContext(context.Background())
	gracefulShutdown(ctx, g, log)

	g.Go(func() error {
		log.Info("starting http server on port", httpPort)
		defer log.Info("closing http server on port", httpPort)

		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err = httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	g.Go(func() error {
		log.Info("starting grpc server on port", grpcPort)
		defer log.Info("closing grpc server on port", grpcPort)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shCtx); err != nil {
				log.Infof("can't close http server listening on %s: %s", httpServer.Addr, err.Error())
			}

			close(errCh)
		}()
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				errCh <- err
			}
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})
	if err := g.Wait(); err != nil {
		log.Infof("gracefully shutting down the servers: %s\n", err.Error())
	}
}

func gracefulShutdown(ctx context.Context, g *errgroup.Group, log logger.Logger) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	g.Go(func() error {
		select {
		case s := <-signals:
			log.Infof("captured signal %v", s)
			return fmt.Errorf("captured signal %v", s)
		case <-ctx.Done():
			return nil
		}
	})
}
