package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/YaroslavGaponov/grotto/internal/configure"
	"github.com/YaroslavGaponov/grotto/internal/masterservice"
	"github.com/YaroslavGaponov/grotto/pkg/logger"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	config := configure.NewConfigure()
	log := logger.NewLogger(config.LogLevel) 

	log.Infof("master service is starting...")
	master := masterservice.New(config.MasterServiceAddr, config.ChunkServiceUrls)
	go master.Start()

	<-sigs

	master.Stop()
}
