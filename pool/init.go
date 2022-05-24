package pool

import (
	"container/list"
	"context"
	"gopkg.in/errgo.v2/fmt/errors"
	"log"
	"transactionSystemTestTask/service"
	"transactionSystemTestTask/store"
)

func Initialisation(ctx context.Context, store *store.Store) error {

	// Get pending transactions
	pending, err := service.NewTransactionWebService(ctx, store).GetPending(ctx)
	if err != nil {
		return errors.Wrap(err)
	}

	log.Println("Pulling pending transactions")

	queues := &Mapping
	for _, transaction := range *pending {
		loaded, ok := queues.Load(transaction.ClientId)
		if !ok {
			queues.Store(transaction.ClientId, list.New())
			loaded, ok = queues.Load(transaction.ClientId)
		}
		queue := loaded.(*list.List)
		queue.PushBack(transaction)
	}

	log.Println("Pool initialised")
	return nil
}
