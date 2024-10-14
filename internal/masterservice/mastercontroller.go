package masterservice

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/YaroslavGaponov/grotto/pkg/common"
	"github.com/YaroslavGaponov/grotto/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	CHUNK_SIZE   = 8 * 1024
	METADATA_TAG = ".metadata"
)

type MasterController struct {
	log                 logger.ILogger
	chunkServiceClients []ChunkServiceClient
	channels            map[string]chan common.Event
	upgrader            websocket.Upgrader
	replicas            int
}

func NewMasterController(log logger.ILogger, chunkServiceUrls []string, replicas int) MasterController {
	chunkServiceClients := make([]ChunkServiceClient, len(chunkServiceUrls))
	for idx, url := range chunkServiceUrls {
		chunkServiceClients[idx] = NewChunkServiceClient(url)
	}
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return MasterController{
		log:                 log,
		chunkServiceClients: chunkServiceClients,
		channels:            make(map[string]chan common.Event),
		upgrader:            upgrader,
		replicas:            replicas,
	}
}

func (masterController *MasterController) Save(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if len(name) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	metadata := NewFileMetadata()

	buf := make([]byte, CHUNK_SIZE)
	id := 0
	for {

		n, err := io.ReadFull(r.Body, buf)

		if err == io.EOF {
			break
		}

		if err != nil && err != io.ErrUnexpectedEOF {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		idx := rand.Intn(len(masterController.chunkServiceClients))
		for replica := 0; replica < masterController.replicas; replica++ {
			client := masterController.chunkServiceClients[idx]
			client.Save(name, id, buf[:n])
			metadata.AddChunk(id, client.Url)
			if idx >= len(masterController.chunkServiceClients)-1 {
				idx = 0
			} else {
				idx++
			}
		}
		id++
	}
	for _, client := range masterController.chunkServiceClients {
		err := client.Save(name+METADATA_TAG, 0, metadata.ToByteArray())
		if err != nil {
			http.Error(w, "Error writing body to chunk service", http.StatusInternalServerError)
			return
		}
	}

	for _, channel := range masterController.channels {
		channel <- common.Event{File: name, Action: common.ACTION_ADD}
	}
}

func (masterController *MasterController) Load(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if len(name) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	idx := rand.Intn(len(masterController.chunkServiceClients))
	client := masterController.chunkServiceClients[idx]
	metadataBody, err := client.Load(name+METADATA_TAG, 0)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	metadata := NewFileMetadata()
	metadata.Load(metadataBody)

	for id := 0; ; id++ {
		urls, found := metadata.chunks[id]
		if !found {
			break
		}
		idx := rand.Intn(len(urls))
		client := NewChunkServiceClient(urls[idx])
		data, err := client.Load(name, id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Write(data)
	}
}

func (masterController *MasterController) List(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]struct{})
	for _, client := range masterController.chunkServiceClients {
		res, err := client.List()
		if err != nil {
			http.Error(w, "Error writing body to chunk service", http.StatusInternalServerError)
			return
		}
		for _, r := range res {
			m[r.Name] = struct{}{}
		}
	}
	var list []string
	for name := range m {
		if !strings.HasSuffix(name, METADATA_TAG) {
			list = append(list, name)
		}
	}
	result, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "Error writing body to chunk service", http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func (masterController *MasterController) Remove(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if len(name) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	idx := rand.Intn(len(masterController.chunkServiceClients))
	client := masterController.chunkServiceClients[idx]
	metadataBody, err := client.Load(name+METADATA_TAG, 0)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	metadata := NewFileMetadata()
	metadata.Load(metadataBody)

	for id := 0; ; id++ {
		urls, found := metadata.chunks[id]
		if !found {
			break
		}
		for _, url := range urls {
			client := NewChunkServiceClient(url)
			err := client.Remove(name, id)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
	}

	for _, client := range masterController.chunkServiceClients {
		err := client.Remove(name+METADATA_TAG, 0)
		if err != nil {
			http.Error(w, "Error writing body to chunk service", http.StatusInternalServerError)
			return
		}
	}

	for _, channel := range masterController.channels {
		channel <- common.Event{File: name, Action: common.ACTION_REMOVE}
	}
}

func (masterController *MasterController) Events(w http.ResponseWriter, r *http.Request) {
	c, err := masterController.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Upgrade error", http.StatusInternalServerError)
		return
	}
	defer c.Close()

	id := uuid.New()
	channel := make(chan common.Event)
	masterController.channels[id.String()] = channel
	defer delete(masterController.channels, id.String())

	for {
		event := <-channel
		if data, err := json.Marshal(event); err == nil {
			c.WriteMessage(websocket.TextMessage, data)
		}
	}
}
