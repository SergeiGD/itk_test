package http

import (
	"context"
	"errors"

	"github.com/SergeiGD/golang-template/internal/usecases"
	"github.com/SergeiGD/golang-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type walletsHandlers struct {
	walletsUseCase usecases.WalletsUseCases
	logger         *logger.Logger
}

func NewWalletHandlers(uc usecases.WalletsUseCases, logger *logger.Logger) StrictServerInterface {
	return &walletsHandlers{walletsUseCase: uc, logger: logger}
}

func (h *walletsHandlers) GetWalletById(ctx context.Context, request GetWalletByIdRequestObject) (GetWalletByIdResponseObject, error) {
	item, err := h.walletsUseCase.GetWalletById(ctx, request.Id)
	if err != nil {

		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).CtxError(ctx, "error on getting wallet by id request")

		if errors.Is(err, usecases.ErrOutOfLimits) {
			return GetWalletById500JSONResponse{
				N50xErrorJSONResponse{Detail: "Превышен лимит запросов"},
			}, nil
		}

		return GetWalletByIddefaultJSONResponse{
			StatusCode: 400,
			Body:       ResponseError{Detail: err.Error()},
		}, nil
	}

	if item == nil {
		return GetWalletById404JSONResponse{
			N404NotFoundJSONResponse{Detail: "Не найдено"},
		}, nil
	}

	return GetWalletById200JSONResponse(mapDomainToPresentationWallet(item)), nil

}

func (h *walletsHandlers) WalletOperation(ctx context.Context, request WalletOperationRequestObject) (WalletOperationResponseObject, error) {
	validate := validator.New()
	err := validate.Struct(request.Body)
	if err != nil {
		return WalletOperationdefaultJSONResponse{
			StatusCode: 400,
			Body:       ResponseError{Detail: err.Error()},
		}, nil
	}
	req := mapMakeWalletOperation(*request.Body)
	item, err := h.walletsUseCase.MakeWalletOperation(ctx, req.Id, req.Operation, req.Amount)

	if err != nil {

		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).CtxError(ctx, "error on making wallet operation request")

		if errors.Is(err, usecases.ErrOutOfLimits) {
			return WalletOperation500JSONResponse{
				N50xErrorJSONResponse{Detail: "Превышен лимит запросов"},
			}, nil
		}

		return WalletOperationdefaultJSONResponse{
			StatusCode: 400,
			Body:       ResponseError{Detail: err.Error()},
		}, nil
	}

	if item == nil {
		return WalletOperation404JSONResponse{
			N404NotFoundJSONResponse{Detail: "Не найдено"},
		}, nil
	}

	return WalletOperation200JSONResponse(mapDomainToPresentationWallet(item)), nil

}
