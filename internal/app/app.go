package app

import (
	"github.com/SergeiGD/itk_test/config"
	"github.com/SergeiGD/itk_test/internal/server"
)

type App struct {
	cfg        *config.Config
	HttpServer server.IServer
}

func NewApp(cfg *config.Config, httpServer server.IServer) *App {
	return &App{
		cfg:        cfg,
		HttpServer: httpServer,
	}
}
