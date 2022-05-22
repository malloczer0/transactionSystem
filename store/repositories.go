package store

import (
	"context"
	"github.com/google/uuid"
	"transactionSystemTestTask/model"
)

type TransactionRepo interface {
	GetTransaction(context.Context, uuid.UUID) (*model.DBTransaction, error)
	GetPending(ctx context.Context) (*[]model.DBTransaction, error)
	CreateTransaction(context.Context, *model.DBTransaction) (*model.DBTransaction, error)
	UpdateTransaction(context.Context, *model.DBTransaction) (*model.DBTransaction, error)
}

type ClientRepo interface {
	GetClient(context.Context, uuid.UUID) (*model.DBClient, error)
	CreateClient(context.Context, *model.DBClient) (*model.DBClient, error)
	UpdateClient(context.Context, *model.DBClient) (*model.DBClient, error)
	DeleteClient(context.Context, uuid.UUID) error
}
