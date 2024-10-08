package masterservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/YaroslavGaponov/grotto/internal/driver"
)

type ChunkServiceClient struct {
	Url string
}

func NewChunkServiceClient(url string) ChunkServiceClient {
	return ChunkServiceClient{
		Url: url,
	}
}

func (client *ChunkServiceClient) Save(name string, id int, data []byte) error {
	url := fmt.Sprintf("%s/v1/chunk/%s/%d", client.Url, name, id)
	resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *ChunkServiceClient) Load(name string, id int) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/chunk/%s/%d", client.Url, name, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (client *ChunkServiceClient) List() ([]driver.ChunkId, error) {
	url := fmt.Sprintf("%s/v1/chunks/all", client.Url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chunks []driver.ChunkId
	if err := json.Unmarshal(data, &chunks); err != nil {
		return nil, err
	}
	return chunks, nil

}

func (client *ChunkServiceClient) Remove(name string, id int) error {
	url := fmt.Sprintf("%s/v1/chunk/%s/%d", client.Url, name, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}
