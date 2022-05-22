package pool

import (
	"container/list"
	"context"
	"github.com/google/uuid"
	"gopkg.in/errgo.v2/fmt/errors"
	"transactionSystemTestTask/service"
	"transactionSystemTestTask/store"
)

func Initialisation(ctx context.Context, store *store.Store) (chan map[uuid.UUID]*list.List, error) {

	// Get pending transactions
	pending, err := service.NewTransactionWebService(ctx, store).GetPending(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var pool map[uuid.UUID]*list.List
	for _, transaction := range *pending {
		if pool[transaction.ID] == nil {
			pool[transaction.ID] = list.New()
		}
		pool[transaction.ID].PushBack(transaction)
	}
	var channel = make(chan map[uuid.UUID]*list.List)
	return channel, nil
}
