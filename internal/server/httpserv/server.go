package httpserv

import (
	"context"

	"github.com/SergeiGD/itk_test/config"
	http2 "github.com/SergeiGD/itk_test/internal/api/http"
	"github.com/SergeiGD/itk_test/internal/di"
	"github.com/SergeiGD/itk_test/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	middleware "github.com/oapi-codegen/gin-middleware"
)

type HttpServer struct {
	cfg    *config.Config
	logger *logger.Logger
}

func NewHttpServer(cfg *config.Config, logger *logger.Logger) *HttpServer {
	return &HttpServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	swagger, err := http2.GetSwagger()
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error on init swagger schema")
		return err
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	// config.AllowOriginFunc = func(origin string) bool {}
	config.AllowOrigins = []string{"http://0.0.0.0:8081"}
	r.Use(cors.New(config))
	r.Use(middleware.OapiRequestValidator(swagger))

	repo, err := di.InitializeWalletRepo(*s.cfg, s.logger)

	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error on init repo")
		return err
	}

	service := di.InitializeWalletService(repo, *s.cfg, s.logger)
	limiter := di.InitializeLimiterService(*s.cfg)

	handler := http2.NewStrictHandler(
		http2.NewWalletHandlers(di.InitializeWalletUseCases(service, limiter, *s.cfg, s.logger), s.logger),
		make([]http2.StrictMiddlewareFunc, 0),
	)

	http2.RegisterHandlers(r, handler)

	err = r.Run(":8080")

	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error on starting gin server")
		return err
	}

	return nil
}
