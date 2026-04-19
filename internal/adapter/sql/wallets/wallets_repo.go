package wallets

import (
	"context"

	"github.com/SergeiGD/itk_test/internal/models"
	"github.com/SergeiGD/itk_test/pkg/logger"
	"github.com/SergeiGD/itk_test/pkg/postgres"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletsRepository interface {
	GetWalletById(ctx context.Context, walletId uuid.UUID) (*models.WalletModel, error)
	WithdrawFromWallet(ctx context.Context, walletId uuid.UUID, amount decimal.Decimal) (*models.WalletModel, error)
	DepositToWallet(ctx context.Context, walletId uuid.UUID, amount decimal.Decimal) (*models.WalletModel, error)
}

type walletsRepository struct {
	client postgres.Client
	logger *logger.Logger
}

func NewWalletsRepository(client postgres.Client, logger *logger.Logger) WalletsRepository {
	return &walletsRepository{
		client: client,
		logger: logger,
	}
}
