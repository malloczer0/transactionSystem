package pool

import (
	"container/list"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sync"
	"transactionSystemTestTask/model"
	globalMutex "transactionSystemTestTask/mutex"
	"transactionSystemTestTask/service"
)

func Handler(ctx context.Context, svc service.Manager, channel chan map[uuid.UUID]*list.List) {
	for {
		var pool = <-channel
		if len(pool) == 0 {
			continue
		}
		for id, queue := range pool {
			var mutex = globalMutex.QueueStateMutex(id)
			mutex.Lock()
			if (queue).Len() == 0 {
				pool[id] = nil
			}
			mutex.Unlock()

			var channel = make(chan *list.List)
			channel <- queue
			go handleTransactionQueue(ctx, svc, id, channel, mutex)
		}
	}
}

func handleTransactionQueue(ctx context.Context, svc service.Manager, uuid uuid.UUID, queue chan *list.List, mutex *sync.Mutex) {
	for {
		mutex.Lock()
		if (<-queue).Len() == 0 {
			globalMutex.PopMutex(uuid)
			return
		}
		transaction := (<-queue).Front()
		mutex.Unlock()
		err := processSQL(ctx, svc, transaction.Value.(*model.Transaction))
		if err != nil {
			return
		}
		mutex.Lock()
		(<-queue).Remove(transaction)
		mutex.Unlock()
	}
}

func Update(ctx context.Context, svc service.Manager, transaction *model.Transaction) (*model.Transaction, error) {
	var mutex = globalMutex.QueueStateMutex(transaction.ClientId)
	mutex.Lock()
	id := &transaction.ClientId
	var queue = <-QueueState(*id)
	queue.PushBack(transaction)
	UpdateQueueState(*id, queue)
	mutex.Unlock()
	createTransaction, err := svc.Transaction.CreateTransaction(ctx, transaction)
	return createTransaction, err
}

func processSQL(ctx context.Context, svc service.Manager, transaction *model.Transaction) error {
	// get client
	client, err := svc.Client.GetClient(ctx, transaction.ClientId)
	if err != nil {
		return errors.Wrap(err, "svc.client.GetClient")
	}
	if client.Balance > transaction.Change {
		transaction.Status = 1
		_, err := svc.Transaction.UpdateTransaction(ctx, transaction)
		if err != nil {
			return errors.Wrap(err, "svc.transaction.UpdateTransaction")
		}

		client.Balance = client.Balance - transaction.Change
		_, err = svc.Client.UpdateClient(ctx, client)
		if err != nil {
			return errors.Wrap(err, "svc.client.UpdateClient")
		}
	} else {
		transaction.Status = -1
		_, err := svc.Transaction.UpdateTransaction(ctx, transaction)
		if err != nil {
			return errors.Wrap(err, "svc.transaction.UpdateTransaction")
		}
	}
	return nil
}
