//go:build wireinject
// +build wireinject

package di

import (
	"github.com/SergeiGD/golang-template/config"
	"github.com/SergeiGD/golang-template/internal/adapter/sql/wallets"
	"github.com/SergeiGD/golang-template/internal/services"
	"github.com/SergeiGD/golang-template/internal/usecases"
	"github.com/SergeiGD/golang-template/pkg/logger"
	"github.com/google/wire"
)

func InitializeWalletRepo(cfg config.Config, logger *logger.Logger) (wallets.WalletsRepository, error) {
	wire.Build(ProvideWalletsRepository)
	return nil, nil
}

func InitializeWalletService(repo wallets.WalletsRepository, cfg config.Config, logger *logger.Logger) services.WalletsService {
	wire.Build(ProvideWalletsService)
	return nil
}

func InitializeLimiterService(cfg config.Config) services.WalletRateLimiter {
	wire.Build(ProvideLimiterService)
	return nil
}

func InitializeWalletUseCases(service services.WalletsService, limiter services.WalletRateLimiter, cfg config.Config, logger *logger.Logger) usecases.WalletsUseCases {
	wire.Build(ProvideWalletsUseCases)
	return nil
}
