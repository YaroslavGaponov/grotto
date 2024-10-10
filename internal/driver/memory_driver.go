package driver

type MemoryDriver struct {
	store map[ChunkId][]byte
}

func NewMemoryDriver() *MemoryDriver {
	return &MemoryDriver{
		store: make(map[ChunkId][]byte),
	}
}

func (driver *MemoryDriver) Save(chunkId ChunkId, data []byte) error {
	driver.store[chunkId] = data
	return nil
}

func (driver *MemoryDriver) Load(chunkId ChunkId) ([]byte, error) {
	if data, found := driver.store[chunkId]; found {
		return data, nil
	}
	return nil, errNotFound
}

func (driver *MemoryDriver) Remove(chunkId ChunkId) error {
	if _, found := driver.store[chunkId]; found {
		delete(driver.store, chunkId)
		return nil
	}
	return errNotFound
}

func (driver *MemoryDriver) List() ([]ChunkId, error) {
	list := make([]ChunkId, 0)
	for id := range driver.store {
		list = append(list, id)
	}
	return list, nil
}
