package masterservice

import (
	"errors"
	"math/rand"
)

var (
	ErrNotFound = errors.New("chunk service client is not found")
)

type ChunkServiceClientPool struct {
	Clients []ChunkServiceClient
	Urls    map[string]int
}

func NewChunkServiceClientPool(urls []string) ChunkServiceClientPool {
	chunkServiceClientPool := ChunkServiceClientPool{
		Clients: make([]ChunkServiceClient, 0, len(urls)),
		Urls:    make(map[string]int, len(urls)),
	}
	for _, url := range urls {
		chunkServiceClientPool.Add(url)
	}
	return chunkServiceClientPool
}

func (c *ChunkServiceClientPool) Add(url string) {
	idx := len(c.Clients)
	client := NewChunkServiceClient(url)
	c.Clients = append(c.Clients, client)
	c.Urls[url] = idx
}

func (c *ChunkServiceClientPool) Get(url string) (*ChunkServiceClient, error) {
	if idx, found := c.Urls[url]; found {
		return &c.Clients[idx], nil
	}
	return nil, ErrNotFound
}

func (c *ChunkServiceClientPool) GetAll() []ChunkServiceClient {
	return c.Clients
}

func (c *ChunkServiceClientPool) GetRandom() (*ChunkServiceClient, error) {
	if len(c.Clients) == 0 {
		return nil, ErrNotFound
	}
	idx := rand.Intn(len(c.Clients))
	return &c.Clients[idx], nil
}
