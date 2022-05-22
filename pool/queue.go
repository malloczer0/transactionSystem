package pool

import (
	"container/list"
	"github.com/google/uuid"
)

var mapping map[uuid.UUID]*list.List

func QueueState(clientId uuid.UUID) chan *list.List {
	channel := make(chan *list.List)
	channel <- mapping[clientId]
	return channel
}

func UpdateQueueState(clientId uuid.UUID, queue *list.List) {
	mapping[clientId] = queue
}
