package services

import (
	"context"
	"errors"

	"github.com/SergeiGD/itk_test/internal/adapter/sql/wallets"
	"github.com/SergeiGD/itk_test/internal/models"
	"github.com/SergeiGD/itk_test/pkg/logger"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
	ErrNegativeBalance  = errors.New("balance cant be negative")
)

type WalletsService interface {
	GetWalletById(ctx context.Context, walletId uuid.UUID) (*models.WalletModel, error)
	MakeWalletOperation(ctx context.Context, walletId uuid.UUID, operation models.OperationType, amount decimal.Decimal) (*models.WalletModel, error)
}

type walletsService struct {
	repo   wallets.WalletsRepository
	logger *logger.Logger
}

func NewWalletsService(repo wallets.WalletsRepository, logger *logger.Logger) WalletsService {
	return &walletsService{
		repo:   repo,
		logger: logger,
	}
}

func (s *walletsService) GetWalletById(ctx context.Context, walletId uuid.UUID) (*models.WalletModel, error) {
	wallet, err := s.repo.GetWalletById(ctx, walletId)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletsService) MakeWalletOperation(ctx context.Context, walletId uuid.UUID, operation models.OperationType, amount decimal.Decimal) (*models.WalletModel, error) {
	if operation == models.WithdrawOperation {
		wallet, err := s.repo.GetWalletById(ctx, walletId)
		if err != nil {
			return nil, err
		}
		if wallet.Balance.Sub(amount).IsNegative() {
			return nil, ErrNegativeBalance
		}
		res, err := s.repo.WithdrawFromWallet(ctx, walletId, amount)
		return res, err
	}
	if operation == models.DepositOperation {
		res, err := s.repo.DepositToWallet(ctx, walletId, amount)
		return res, err
	}
	return nil, ErrInvalidOperation
}
