package driver

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DiskDriver struct {
	root string
}

func NewDiskDriver(root string) *DiskDriver {
	return &DiskDriver{
		root: strings.TrimSuffix(root, "/"),
	}
}

func (driver *DiskDriver) Save(chunkId ChunkId, data []byte) error {
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
	name := fmt.Sprintf("%s/%s/%d", driver.root, chunkId.Name, chunkId.Id)
	return os.ReadFile(name)
}

func (driver *DiskDriver) Remove(chunkId ChunkId) error {
	name := fmt.Sprintf("%s/%s/%d", driver.root, chunkId.Name, chunkId.Id)
	return os.Remove(name)
}

func (driver *DiskDriver) List() ([]ChunkId, error) {
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
		if err := os.Chdir(name); err != nil {
			return err
		}
	}
	return nil
}
