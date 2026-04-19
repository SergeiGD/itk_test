package http

import (
	"github.com/SergeiGD/itk_test/internal/dto"
	"github.com/SergeiGD/itk_test/internal/models"
	"github.com/shopspring/decimal"
)

func mapDomainToPresentationWallet(item *models.WalletModel) Wallet {
	balance, _ := item.Balance.Float64()
	return Wallet{
		Balance:   float32(balance),
		CreatedAt: item.CreatedAt,
	}
}

func mapMakeWalletOperation(req WalletOperationJSONRequestBody) dto.MakeWalletOperationDTO {
	return dto.MakeWalletOperationDTO{
		Id:        req.WalletId,
		Operation: models.OperationType(req.OperationType),
		Amount:    decimal.NewFromFloat(float64(req.Amount)),
	}
}
