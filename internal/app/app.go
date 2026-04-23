package app

import (
	"github.com/SergeiGD/golang-template/config"
	"github.com/SergeiGD/golang-template/internal/server"
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
