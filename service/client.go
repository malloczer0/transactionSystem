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

// ClientWebService ...
type ClientWebService struct {
	ctx   context.Context
	store *store.Store
}

// NewClientWebService creates a new transaction web service
func NewClientWebService(ctx context.Context, store *store.Store) *ClientWebService {
	return &ClientWebService{
		ctx:   ctx,
		store: store,
	}
}

// GetClient ...
func (clientWebService *ClientWebService) GetClient(ctx context.Context, id uuid.UUID) (*model.Client, error) {
	db, err := clientWebService.store.Client.GetClient(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "clientWebService.client.GetClient")
	}
	if db == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("Client '%s' not found", id.String()))
	}

	return db.ToWeb(), nil
}

// CreateClient ...
func (clientWebService ClientWebService) CreateClient(ctx context.Context, client *model.Client) (*model.Client, error) {
	client.ID = uuid.New()

	_, err := clientWebService.store.Client.CreateClient(ctx, client.ToDB())
	if err != nil {
		return nil, errors.Wrap(err, "clientWebService.client.CreateClient error")
	}

	// get created by ID
	createdDBClient, err := clientWebService.store.Client.GetClient(ctx, client.ID)
	if err != nil {
		return nil, errors.Wrap(err, "clientWebService.client.GetClient error")
	}

	return createdDBClient.ToWeb(), nil
}

// UpdateClient ...
func (clientWebService *ClientWebService) UpdateClient(ctx context.Context, client *model.Client) (*model.Client, error) {
	clientDB, err := clientWebService.store.Client.GetClient(ctx, client.ID)
	if err != nil {
		return nil, errors.Wrap(err, "clientWebService.clientDB.GetClient error")
	}
	if clientDB == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("Client '%s' not found", client.ID.String()))
	}

	// update
	_, err = clientWebService.store.Client.UpdateClient(ctx, client.ToDB())
	if err != nil {
		return nil, errors.Wrap(err, "clientWebService.client.UpdateClient error")
	}

	// get updated by ID
	updatedDBClient, err := clientWebService.store.Client.GetClient(ctx, client.ID)
	if err != nil {
		return nil, errors.Wrap(err, "clientWebService.client.GetClient error")
	}

	return updatedDBClient.ToWeb(), nil
}
