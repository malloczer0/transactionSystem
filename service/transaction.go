package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"transactionSystemTestTask/lib/types"
	"transactionSystemTestTask/model"
	"transactionSystemTestTask/store"
)

// TransactionWebService ...
type TransactionWebService struct {
	ctx   context.Context
	store *store.Store
}

// NewTransactionWebService creates a new transaction web service
func NewTransactionWebService(ctx context.Context, store *store.Store) *TransactionWebService {
	return &TransactionWebService{
		ctx:   ctx,
		store: store,
	}
}

// GetTransaction ...
func (transactionWebService *TransactionWebService) GetTransaction(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	db, err := transactionWebService.store.Transaction.GetTransaction(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "transactionWebService.transaction.GetTransaction")
	}
	if db == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("Transaction '%s' not found", id.String()))
	}

	return db.ToWeb(), nil
}

// GetPending ...
func (transactionWebService *TransactionWebService) GetPending(ctx context.Context) (*[]model.Transaction, error) {
	db, err := transactionWebService.store.Transaction.GetPending(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "transactionWebService.transaction.GetPending")
	}
	if db == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("Pending transactions not found"))
	}

	var transactions []model.Transaction
	for index, transaction := range *db {
		transactions[index] = *transaction.ToWeb()
	}

	return &transactions, nil
}

// CreateTransaction ...
func (transactionWebService TransactionWebService) CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	transaction.ID = uuid.New()
	transaction.Status = 0

	_, err := transactionWebService.store.Transaction.CreateTransaction(ctx, transaction.ToDB())
	if err != nil {
		return nil, errors.Wrap(err, "transactionWebService.transaction.CreateTransaction error")
	}

	// get created by ID
	createdDBTransaction, err := transactionWebService.store.Transaction.GetTransaction(ctx, transaction.ID)
	if err != nil {
		return nil, errors.Wrap(err, "transactionWebService.transaction.GetTransaction error")
	}

	return createdDBTransaction.ToWeb(), nil
}

// UpdateTransaction ...
func (transactionWebService *TransactionWebService) UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	transactionDB, err := transactionWebService.store.Transaction.GetTransaction(ctx, transaction.ID)
	if err != nil {
		return nil, errors.Wrap(err, "transactionWebService.transaction.GetTransaction error")
	}
	if transactionDB == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("Client '%s' not found", transaction.ID.String()))
	}

	// update
	_, err = transactionWebService.store.Transaction.UpdateTransaction(ctx, transaction.ToDB())
	if err != nil {
		return nil, errors.Wrap(err, "transactionWebService.transaction.UpdateClient error")
	}

	// get updated by ID
	updatedDBTransaction, err := transactionWebService.store.Transaction.GetTransaction(ctx, transaction.ID)
	if err != nil {
		return nil, errors.Wrap(err, "transactionWebService.transaction.GetTransaction error")
	}

	return updatedDBTransaction.ToWeb(), nil
}
