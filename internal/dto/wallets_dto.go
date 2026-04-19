package dto

import (
	"github.com/SergeiGD/itk_test/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MakeWalletOperationDTO struct {
	Id        uuid.UUID
	Amount    decimal.Decimal
	Operation models.OperationType
}
