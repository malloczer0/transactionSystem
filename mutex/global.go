package mutex

import (
	"github.com/google/uuid"
	"sync"
)

var mapping map[uuid.UUID]*sync.Mutex

func QueueStateMutex(clientId uuid.UUID) *sync.Mutex {
	if mapping[clientId] == nil {
		var m sync.Mutex
		mapping[clientId] = &m
	}
	return mapping[clientId]
}

func PopMutex(uuid uuid.UUID) {
	mapping[uuid] = nil
}
