package driver

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type DiskDriver struct {
	mu   sync.RWMutex
	root string
}

func NewDiskDriver(root string) *DiskDriver {
	return &DiskDriver{
		root: strings.TrimSuffix(root, "/"),
	}
}

func (driver *DiskDriver) Save(chunkId ChunkId, data []byte) error {
	driver.mu.Lock()
	defer driver.mu.Unlock()

	if err := createDir(driver.root); err != nil {
		return err
	}

	name := fmt.Sprintf("%s/%s", driver.root, chunkId.Name)
	if err := createDir(name); err != nil {
		return err
	}
	if err := os.WriteFile(fmt.Sprintf("%s/%d", name, chunkId.Id), data, 0644); err != nil {
		return err
	}
	return nil
}

func (driver *DiskDriver) Load(chunkId ChunkId) ([]byte, error) {
	driver.mu.RLock()
	defer driver.mu.RUnlock()

	name := fmt.Sprintf("%s/%s/%d", driver.root, chunkId.Name, chunkId.Id)
	return os.ReadFile(name)
}

func (driver *DiskDriver) Remove(chunkId ChunkId) error {
	driver.mu.Lock()
	defer driver.mu.Unlock()

	dir := fmt.Sprintf("%s/%s", driver.root, chunkId.Name)
	name := fmt.Sprintf("%s/%d", dir, chunkId.Id)

	if err := os.Remove(name); err != nil {
		return err
	}
	list, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return os.Remove(dir)
	}
	return nil
}

func (driver *DiskDriver) List() ([]ChunkId, error) {
	driver.mu.RLock()
	defer driver.mu.RUnlock()

	files, err := os.ReadDir(driver.root)
	if err != nil {
		return nil, err
	}
	var list []ChunkId
	for _, file := range files {
		chunks, err := os.ReadDir(driver.root + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		for _, chunk := range chunks {
			id, err := strconv.Atoi(chunk.Name())
			if err != nil {
				return nil, err
			}
			list = append(list, ChunkId{file.Name(), id})
		}
	}
	return list, nil
}

func createDir(name string) error {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		if err := os.Mkdir(name, 0644); err != nil {
			return err
		}
	}
	return nil
}
