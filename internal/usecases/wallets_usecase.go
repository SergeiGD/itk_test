package usecases

import (
	"context"
	"errors"

	"github.com/SergeiGD/golang-template/internal/models"
	"github.com/SergeiGD/golang-template/internal/services"
	"github.com/SergeiGD/golang-template/pkg/logger"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var ErrOutOfLimits = errors.New("out of limits")

type WalletsUseCases interface {
	GetWalletById(ctx context.Context, walletId uuid.UUID) (*models.WalletModel, error)
	MakeWalletOperation(ctx context.Context, walletId uuid.UUID, operation models.OperationType, amount decimal.Decimal) (*models.WalletModel, error)
}

type walletsUseCases struct {
	walletsService services.WalletsService
	limiter        services.WalletRateLimiter
	logger         *logger.Logger
}

func NewWalletsUseCases(walletsService services.WalletsService, limiter services.WalletRateLimiter, logger *logger.Logger) WalletsUseCases {
	return &walletsUseCases{
		walletsService: walletsService,
		limiter:        limiter,
		logger:         logger,
	}
}

func (uc *walletsUseCases) GetWalletById(ctx context.Context, walletId uuid.UUID) (*models.WalletModel, error) {
	if !uc.limiter.IsAllowed(walletId) {
		return nil, ErrOutOfLimits
	}
	wallet, err := uc.walletsService.GetWalletById(ctx, walletId)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (uc *walletsUseCases) MakeWalletOperation(ctx context.Context, walletId uuid.UUID, operation models.OperationType, amount decimal.Decimal) (*models.WalletModel, error) {
	if !uc.limiter.IsAllowed(walletId) {
		return nil, ErrOutOfLimits
	}
	wallet, err := uc.walletsService.MakeWalletOperation(ctx, walletId, operation, amount)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
