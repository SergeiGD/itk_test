package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SergeiGD/golang-template/config"
	"github.com/SergeiGD/golang-template/internal/app"
	"github.com/SergeiGD/golang-template/internal/server/httpserv"
	"github.com/SergeiGD/golang-template/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var cfg config.Config
	err := cleanenv.ReadConfig("config/config.yaml", &cfg)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	logger := logger.NewLogger(&cfg)

	app := app.NewApp(
		&cfg,
		httpserv.NewHttpServer(&cfg, logger),
	)
	eg := errgroup.Group{}

	eg.Go(func() error { return app.HttpServer.Run(ctx) })

	return eg.Wait()

}
