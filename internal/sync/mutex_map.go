package sync

import (
	"sync"

	"github.com/pkg/errors"
)

type MutexMap struct {
	sync.Mutex
	entries map[string]*mutex
}

type mutex struct {
	parent *MutexMap
	sync.Mutex
	references int
	key        string
}

func NewMutexMap() *MutexMap {
	return &MutexMap{
		entries: make(map[string]*mutex),
	}
}

func (m *MutexMap) Entry(key string) sync.Locker {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	entry, ok := m.entries[key]
	if !ok {
		entry = &mutex{
			parent: m,
			key:    key,
		}
		m.entries[key] = entry
	}

	entry.references++
	return entry
}

func (e *mutex) Unlock() {
	parent := e.parent
	parent.Mutex.Lock()

	entry, ok := parent.entries[e.key]
	if !ok {
		panic(errors.Errorf("unlock called on non-existing key '%s'", e.key))
	}
	entry.references--
	if entry.references < 1 {
		delete(parent.entries, e.key)
	}
	parent.Mutex.Unlock()
	e.Mutex.Unlock()
}
