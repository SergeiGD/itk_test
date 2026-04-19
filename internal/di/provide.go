package di

import (
	"context"

	"github.com/SergeiGD/itk_test/config"
	"github.com/SergeiGD/itk_test/internal/adapter/sql/wallets"
	"github.com/SergeiGD/itk_test/internal/services"
	"github.com/SergeiGD/itk_test/internal/usecases"
	"github.com/SergeiGD/itk_test/pkg/logger"
	"github.com/SergeiGD/itk_test/pkg/postgres"
)

func ProvideWalletsRepository(cfg config.Config, logger *logger.Logger) (wallets.WalletsRepository, error) {
	postgresClient, err := postgres.NewClient(
		context.Background(),
		cfg,
	)

	if err != nil {
		return nil, err
	}

	return wallets.NewWalletsRepository(postgresClient, logger), nil
}

func ProvideWalletsService(walletsRepo wallets.WalletsRepository, logger *logger.Logger) services.WalletsService {
	return services.NewWalletsService(walletsRepo, logger)
}

func ProvideLimiterService(cfg config.Config) services.WalletRateLimiter {
	return services.NewWalletRateLimiter(cfg)
}

func ProvideWalletsUseCases(walletsService services.WalletsService, limiter services.WalletRateLimiter, logger *logger.Logger) usecases.WalletsUseCases {
	return usecases.NewWalletsUseCases(walletsService, limiter, logger)
}
