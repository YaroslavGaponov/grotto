package driver

import (
	"sync"
)

type MemoryDriver struct {
	store map[ChunkId][]byte
	mu sync.RWMutex
}

func NewMemoryDriver() *MemoryDriver {
	return &MemoryDriver{
		store: make(map[ChunkId][]byte),
		mu: sync.RWMutex{},
	}
}

func (driver *MemoryDriver) Save(chunkId ChunkId, data []byte) error {
	driver.mu.Lock()
	defer driver.mu.Unlock()
	driver.store[chunkId] = data
	return nil
}

func (driver *MemoryDriver) Load(chunkId ChunkId) ([]byte, error) {
	driver.mu.RLock()
	defer driver.mu.RUnlock()
	if data, found := driver.store[chunkId]; found {
		return data, nil
	}
	return nil, errNotFound
}

func (driver *MemoryDriver) Remove(chunkId ChunkId) error {
	driver.mu.Lock()
	defer driver.mu.Unlock()
	if _, found := driver.store[chunkId]; found {
		delete(driver.store, chunkId)
		return nil
	}
	return errNotFound
}

func (driver *MemoryDriver) List() ([]ChunkId, error) {
	driver.mu.RLock()
	defer driver.mu.RUnlock()
	list := make([]ChunkId, 0)
	for id := range driver.store {
		list = append(list, id)
	}
	return list, nil
}
