package chunkservice

import (
	"net/http"

	"github.com/go-chi/chi"
)

const (
	CHUNKS_BASE_URL = "/v1/chunks"
	CHUNK_BASE_URL  = "/v1/chunk"
)

type ChunkService struct {
	addr   string
	router *chi.Mux
	chunk  ChunkController
}

func New(addr string) ChunkService {
	chunkService := ChunkService{
		addr:   addr,
		router: chi.NewRouter(),
		chunk:  NewChunkController(),
	}
	chunkService.router.Get(CHUNKS_BASE_URL+"/all", chunkService.chunk.List)
	chunkService.router.Get(CHUNK_BASE_URL+"/{name}/{id}", chunkService.chunk.Load)
	chunkService.router.Post(CHUNK_BASE_URL+"/{name}/{id}", chunkService.chunk.Save)
	chunkService.router.Delete(CHUNK_BASE_URL+"/{name}/{id}", chunkService.chunk.Remove)
	return chunkService
}

func (chunkService *ChunkService) Start() error {
	return http.ListenAndServe(chunkService.addr, chunkService.router)
}

func (chunkService *ChunkService) Stop() {
}
