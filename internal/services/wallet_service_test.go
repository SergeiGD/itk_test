package services

import (
	"context"
	"testing"
	"time"

	"github.com/SergeiGD/golang-template/internal/adapter/sql/wallets"
	"github.com/SergeiGD/golang-template/internal/models"
	"github.com/SergeiGD/golang-template/pkg/logger"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletsMakeOperationService(t *testing.T) {
	repo := wallets.NewMockWalletsRepository(t)
	l := logger.NewSilentTestsLogger()
	service := NewWalletsService(repo, l)

	type args struct {
		walletId  uuid.UUID
		operation models.OperationType
		amount    decimal.Decimal
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
		dbRow   *models.WalletModel
	}{
		{
			name: "Withdraw operation",
			args: args{
				walletId:  parseStringUuid("69359037-9599-48e7-b8f2-48393c019135"),
				operation: models.WithdrawOperation,
				amount:    decimal.NewFromInt(123),
			},
			dbRow: &models.WalletModel{
				Id:        parseStringUuid("69359037-9599-48e7-b8f2-48393c019135"),
				Balance:   decimal.NewFromInt(123),
				CreatedAt: time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "Unknown operation",
			args: args{
				walletId:  parseStringUuid("69359037-9599-48e7-b8f2-48393c019135"),
				operation: "test",
				amount:    decimal.NewFromInt(123),
			},
			dbRow:   nil,
			wantErr: ErrInvalidOperation,
		},
		{
			name: "Negative balance",
			args: args{
				walletId:  parseStringUuid("69359037-9599-48e7-b8f2-48393c019135"),
				operation: models.WithdrawOperation,
				amount:    decimal.NewFromInt(1000),
			},
			dbRow: &models.WalletModel{
				Id:        parseStringUuid("69359037-9599-48e7-b8f2-48393c019135"),
				Balance:   decimal.NewFromInt(500),
				CreatedAt: time.Time{},
			},
			wantErr: ErrNegativeBalance,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repo.On("WithdrawFromWallet", mock.Anything, test.args.walletId, test.args.amount).
				Maybe().
				Return(test.dbRow, nil)
			repo.On("GetWalletById", mock.Anything, test.args.walletId).
				Maybe().
				Return(test.dbRow, nil)

			_, err := service.MakeWalletOperation(context.Background(), test.args.walletId, test.args.operation, test.args.amount)
			assert.ErrorIs(t, test.wantErr, err)

		},
		)
	}

}

func parseStringUuid(val string) uuid.UUID {
	res, _ := uuid.Parse(val)
	return res
}
