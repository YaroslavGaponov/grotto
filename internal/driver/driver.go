package driver

type ChunkId struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Driver interface {
	Save(chunkId ChunkId, data []byte) error
	Load(chunkId ChunkId) ([]byte, error)
	Remove(chunkId ChunkId) error
	List() ([]ChunkId, error)
}
