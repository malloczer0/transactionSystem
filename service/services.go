package service

import (
	"context"
	"github.com/google/uuid"
	"transactionSystemTestTask/model"
)

type ClientService interface {
	GetClient(context.Context, uuid.UUID) (*model.Client, error)
	CreateClient(context.Context, *model.Client) (*model.Client, error)
	UpdateClient(context.Context, *model.Client) (*model.Client, error)
}

type TransactionService interface {
	GetTransaction(context.Context, uuid.UUID) (*model.Transaction, error)
	GetPending(ctx context.Context) (*[]model.Transaction, error)
	CreateTransaction(context.Context, *model.Transaction) (*model.Transaction, error)
	UpdateTransaction(context.Context, *model.Transaction) (*model.Transaction, error)
}
