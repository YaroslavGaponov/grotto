package chunkservice

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/YaroslavGaponov/grotto/internal/configure"
	"github.com/YaroslavGaponov/grotto/internal/driver"
	"github.com/go-chi/chi"
)

type ChunkController struct {
	store driver.Driver
}

func NewChunkController(conf configure.Configure) ChunkController {
	var store driver.Driver

	switch conf.DriverType {
	case "disk":
		store = driver.NewDiskDriver(conf.DiskDriverRootDir)
	case "memory":
		store = driver.NewMemoryDriver()
	}
	return ChunkController{
		store: store,
	}
}

func (chunkController *ChunkController) List(w http.ResponseWriter, r *http.Request) {
	catalog, err := chunkController.store.List()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(catalog); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (chunkController *ChunkController) Load(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if len(name) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	chunkId := driver.ChunkId{Name: name, Id: id}
	data, err := chunkController.store.Load(chunkId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Write(data)
}

func (chunkController *ChunkController) Save(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if len(name) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	chunkId := driver.ChunkId{Name: name, Id: id}

	if err := chunkController.store.Save(chunkId, body); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (chunkController *ChunkController) Remove(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if len(name) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	chunkId := driver.ChunkId{Name: name, Id: id}
	if err := chunkController.store.Remove(chunkId); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
