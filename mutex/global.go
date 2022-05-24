package mutex

import (
	"github.com/google/uuid"
	"sync"
)

var mapping map[uuid.UUID]*sync.Mutex

func GetMutex(id uuid.UUID) *sync.Mutex {
	if mapping == nil {
		mapping = make(map[uuid.UUID]*sync.Mutex)
	}
	if mapping[id] == nil {
		var m sync.Mutex
		mapping[id] = &m
	}
	return mapping[id]
}

func PopMutex(uuid uuid.UUID) {
	mapping[uuid] = nil
}
