package masterservice

import (
	"net/http"

	"github.com/go-chi/chi"
)

const (
	FILES_BASE_URL = "/v1/files"
	FILE_BASE_URL = "/v1/file"
)

type MasterService struct {
	addr string
	router *chi.Mux
	masterController MasterController
}

func New(addr string, chunkServiceUrls []string) MasterService {
	masterService := MasterService{ 
		addr: addr,
		router: chi.NewRouter(),
		masterController: NewMasterController(chunkServiceUrls),
	}
	masterService.router.Post(FILE_BASE_URL+"/{name}", masterService.masterController.Save)
	masterService.router.Get(FILE_BASE_URL+"/{name}", masterService.masterController.Load)
	masterService.router.Delete(FILE_BASE_URL+"/{name}", masterService.masterController.Remove)
	masterService.router.Get(FILES_BASE_URL+"/all", masterService.masterController.List)
	return masterService
}

func (masterService *MasterService) Start() error {
	return http.ListenAndServe(masterService.addr, masterService.router)
}

func (masterService *MasterService) Stop() error {
	return nil
}