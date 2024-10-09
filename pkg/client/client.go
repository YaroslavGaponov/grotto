package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/YaroslavGaponov/grotto/pkg/common"
	"github.com/gorilla/websocket"
)

type Client struct {
	Url string
}

func NewClient(url string) Client {
	return Client{
		Url: url,
	}
}

func (client *Client) Save(name string, data []byte) error {
	url := fmt.Sprintf("%s/v1/file/%s", client.Url, name)
	resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) Load(name string) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/file/%s", client.Url, name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (client *Client) List() ([]string, error) {
	url := fmt.Sprintf("%s/v1/files/all", client.Url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []string
	if err := json.Unmarshal(data, &files); err != nil {
		return nil, err
	}
	return files, nil

}

func (client *Client) Remove(name string) error {
	url := fmt.Sprintf("%s/v1/file/%s", client.Url, name)
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

func (client *Client) Watch(channel chan common.Event) {
	u := url.URL{Scheme: "ws", Host: strings.TrimPrefix(client.Url, "http://"), Path: "/events"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		var event common.Event
		if err := json.Unmarshal(message, &event); err != nil {
			fmt.Println(err)
			return
		}
		channel <- event
	}
}
