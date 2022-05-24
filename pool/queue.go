package pool

import (
	"container/list"
	"github.com/google/uuid"
	"sync"
)

var Mapping sync.Map

func QueueState(clientId uuid.UUID) *list.List {
	loaded, ok := Mapping.Load(clientId)
	if !ok {
		Mapping.Store(clientId, list.New())
		loaded, ok = Mapping.Load(clientId)
	}
	return loaded.(*list.List)
}

func UpdateQueuePoolState(clientId uuid.UUID, queue *list.List) {
	Mapping.Store(clientId, queue)
}
