package http

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/SergeiGD/itk_test/config"
	"github.com/SergeiGD/itk_test/internal/adapter/sql/wallets"
	"github.com/SergeiGD/itk_test/internal/models"
	"github.com/SergeiGD/itk_test/internal/services"
	"github.com/SergeiGD/itk_test/internal/usecases"
	"github.com/SergeiGD/itk_test/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"net/http/httptest"
)

func TestGetWalletByIdHandler(t *testing.T) {
	repo := wallets.NewMockWalletsRepository(t)
	l := logger.NewSilentTestsLogger()
	service := services.NewWalletsService(repo, l)
	limiter := services.NewWalletRateLimiter(config.Config{Limiter: struct {
		MaxLimit      int           `env:"LIMITER_MAX_LIMIT"`
		Burst         int           `env:"LIMITER_BURST"`
		CleanInterval time.Duration `env:"LIMITER_CLEAN_INTERVAL"`
	}{MaxLimit: 100, Burst: 100, CleanInterval: 100}})
	uc := usecases.NewWalletsUseCases(service, limiter, l)

	r := gin.Default()
	handler := NewStrictHandler(
		NewWalletHandlers(uc, logger.NewSilentTestsLogger()),
		make([]StrictMiddlewareFunc, 0),
	)
	RegisterHandlers(r, handler)

	type args struct {
		walletId uuid.UUID
	}

	tests := []struct {
		name   string
		args   args
		result *Wallet
		dbRow  *models.WalletModel
		code   int
	}{
		{
			name: "Not found test",
			args: args{
				walletId: parseStringUuid("1d8a7307-b5f8-4686-b9dc-b752430abbd8"),
			},
			result: nil,
			dbRow:  nil,
			code:   404,
		},
		{
			name: "Item found test",
			args: args{
				walletId: parseStringUuid("69359037-9599-48e7-b8f2-48393c019135"),
			},
			result: &Wallet{
				Balance:   123,
				CreatedAt: time.Time{},
			},
			dbRow: &models.WalletModel{
				Id:        parseStringUuid("69359037-9599-48e7-b8f2-48393c019135"),
				Balance:   decimal.NewFromInt(123),
				CreatedAt: time.Time{},
			},
			code: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repo.On("GetWalletById", mock.Anything, test.args.walletId).
				Return(test.dbRow, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/wallets/%s/", test.args.walletId), nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, test.code)

			var result *Wallet
			err := json.NewDecoder(rec.Body).Decode(&result)
			if test.code == 404 {
				result = nil
			}
			assert.NoError(t, err)
			assert.Equal(t, test.result, result)

		},
		)
	}

}

func parseStringUuid(val string) uuid.UUID {
	res, _ := uuid.Parse(val)
	return res
}
