package pool

import (
	"container/list"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"sync"
	"transactionSystemTestTask/model"
	globalMutex "transactionSystemTestTask/mutex"
	"transactionSystemTestTask/service"
	"transactionSystemTestTask/store"
)

func Handler(ctx context.Context, store *store.Store) {
	log.Println("Running handler")
	for {
		Mapping.Range(func(key, value any) bool {
			log.Println(key, value)

			var (
				id    = key.(uuid.UUID)
				queue = value.(*list.List)
			)

			clientMutex := globalMutex.GetMutex(id)
			clientMutex.TryLock()

			// Check if queue is empty
			log.Println("Queue len is ", queue.Len())
			if queue.Len() == 0 {
				Mapping.Delete(key)
			}

			err := handleTransactionQueue(ctx, store, id, queue, clientMutex)
			if err != nil {
				panic(err)
			}

			clientMutex.Unlock()

			return true
		})
	}
}

func handleTransactionQueue(ctx context.Context, store *store.Store, uuid uuid.UUID, queue *list.List, mutex *sync.Mutex) error {
	for {
		mutex.TryLock()
		if queue.Len() == 0 {
			globalMutex.PopMutex(uuid)
			return nil
		}
		front := queue.Front()
		transaction := front.Value.(model.Transaction)
		mutex.Unlock()

		err := processSQL(ctx, store, &transaction)

		if err != nil {
			return errors.Wrap(err, "pool.ProcessSQL failed")
		}

		mutex.TryLock()
		queue.Remove(front)
		mutex.Unlock()
	}
}

func UpdatePool(ctx context.Context, svc service.Manager, transaction *model.Transaction) (*model.Transaction, error) {
	var mutex = globalMutex.GetMutex(transaction.ClientId)
	mutex.TryLock()
	id := &transaction.ClientId
	queue := QueueState(*id)
	queue.PushBack(transaction)
	UpdateQueuePoolState(*id, queue)
	mutex.Unlock()
	createTransaction, err := svc.Transaction.CreateTransaction(ctx, transaction)
	return createTransaction, err
}

func processSQL(ctx context.Context, store *store.Store, transaction *model.Transaction) error {

	log.Println("Processing transaction:", transaction.ID, "with client id", transaction.ClientId)

	clientService := service.NewClientWebService(ctx, store)
	transactionService := service.NewTransactionWebService(ctx, store)

	// Get client
	client, err := clientService.GetClient(ctx, transaction.ClientId)
	if err != nil {
		return errors.Wrap(err, "svc.client.GetClient")
	}

	if client.Balance > transaction.Change {
		transaction.Status = 1
		_, err := transactionService.UpdateTransaction(ctx, transaction)
		if err != nil {
			return errors.Wrap(err, "svc.transaction.UpdateTransaction")
		}

		client.Balance = client.Balance - transaction.Change
		_, err = clientService.UpdateClient(ctx, client)
		if err != nil {
			return errors.Wrap(err, "svc.client.UpdateClient")
		}
	} else {
		transaction.Status = -1
		_, err := transactionService.UpdateTransaction(ctx, transaction)
		if err != nil {
			return errors.Wrap(err, "svc.transaction.UpdateTransaction")
		}
	}

	return nil
}
