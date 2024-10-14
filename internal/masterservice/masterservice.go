package masterservice

import (
	"net/http"

	"github.com/YaroslavGaponov/grotto/internal/configure"
	"github.com/YaroslavGaponov/grotto/pkg/logger"
	"github.com/go-chi/chi"
)

const (
	FILES_BASE_URL  = "/v1/files"
	FILE_BASE_URL   = "/v1/file"
	EVENTS_BASE_URL = "/events"
)

type MasterService struct {
	log              logger.ILogger
	addr             string
	router           *chi.Mux
	masterController MasterController
}

func New(log logger.ILogger,  conf configure.Configure) MasterService {
	masterService := MasterService{
		log:              log,
		addr:             conf.MasterServiceAddr,
		router:           chi.NewRouter(),
		masterController: NewMasterController(log, conf.ChunkServiceUrls, conf.ChunkReplicas),
	}
	masterService.router.Post(FILE_BASE_URL+"/{name}", masterService.masterController.Save)
	masterService.router.Get(FILE_BASE_URL+"/{name}", masterService.masterController.Load)
	masterService.router.Delete(FILE_BASE_URL+"/{name}", masterService.masterController.Remove)
	masterService.router.Get(FILES_BASE_URL+"/all", masterService.masterController.List)

	masterService.router.Get(EVENTS_BASE_URL, masterService.masterController.Events)

	return masterService
}

func (masterService *MasterService) Start() error {
	masterService.log.Infof("server is starting at %s", masterService.addr)
	return http.ListenAndServe(masterService.addr, masterService.router)
}

func (masterService *MasterService) Stop() error {
	return nil
}
