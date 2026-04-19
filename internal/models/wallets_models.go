package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OperationType string

var (
	DepositOperation  OperationType = "deposit"
	WithdrawOperation OperationType = "withdraw"
)

type WalletModel struct {
	Id        uuid.UUID
	Balance   decimal.Decimal
	CreatedAt time.Time
}
