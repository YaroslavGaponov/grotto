package chunkservice

import (
	"net/http"

	"github.com/YaroslavGaponov/grotto/internal/configure"
	"github.com/YaroslavGaponov/grotto/pkg/logger"
	"github.com/go-chi/chi"
)

const (
	CHUNKS_BASE_URL = "/v1/chunks"
	CHUNK_BASE_URL  = "/v1/chunk"
)

type ChunkService struct {
	log    logger.ILogger
	addr   string
	router *chi.Mux
	chunk  ChunkController
}

func New(log logger.ILogger, conf configure.Configure) ChunkService {
	chunkService := ChunkService{
		log:    log,
		addr:   conf.ChunkServiceAddr,
		router: chi.NewRouter(),
		chunk:  NewChunkController(log, conf),
	}
	chunkService.router.Get(CHUNKS_BASE_URL+"/all", chunkService.chunk.List)
	chunkService.router.Get(CHUNK_BASE_URL+"/{name}/{id}", chunkService.chunk.Load)
	chunkService.router.Post(CHUNK_BASE_URL+"/{name}/{id}", chunkService.chunk.Save)
	chunkService.router.Delete(CHUNK_BASE_URL+"/{name}/{id}", chunkService.chunk.Remove)
	return chunkService
}

func (chunkService *ChunkService) Start() error {
	chunkService.log.Infof("server is starting at %s", chunkService.addr)
	return http.ListenAndServe(chunkService.addr, chunkService.router)
}

func (chunkService *ChunkService) Stop() {
	chunkService.log.Infof("server is stopping")
}
