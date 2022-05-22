package service

import (
	"context"
	"github.com/pkg/errors"
	"transactionSystemTestTask/store"
)

// Manager is just a collection of all services we have in the project
type Manager struct {
	Transaction TransactionService
	Client      ClientService
}

// NewManager creates new service manager
func NewManager(ctx context.Context, store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("No store provided")
	}
	return &Manager{
		Transaction: NewTransactionWebService(ctx, store),
	}, nil
}
