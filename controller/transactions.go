package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"transactionSystemTestTask/lib/types"
	"transactionSystemTestTask/logger"
	"transactionSystemTestTask/model"
	"transactionSystemTestTask/pool"
	"transactionSystemTestTask/service"
)

// TransactionController ...
type TransactionController struct {
	ctx      context.Context
	services *service.Manager
	logger   *logger.Logger
}

// NewTransactions creates a new transaction controller.
func NewTransactions(ctx context.Context, services *service.Manager, logger *logger.Logger) *TransactionController {
	return &TransactionController{
		ctx:      ctx,
		services: services,
		logger:   logger,
	}
}

// Create creates new transaction
func (ctr *TransactionController) Create(ctx echo.Context) error {
	var transaction model.Transaction
	err := ctx.Bind(&transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode transaction data"))
	}
	err = ctx.Validate(&transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	createdTransaction, err := ctr.services.Transaction.CreateTransaction(ctx.Request().Context(), &transaction)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create transaction"))
		}
	}
	pool.UpdatePool(createdTransaction)
	ctr.logger.Debug().Msgf("Created transaction '%s'", createdTransaction.ID.String())
	return ctx.JSON(http.StatusCreated, createdTransaction)
}
