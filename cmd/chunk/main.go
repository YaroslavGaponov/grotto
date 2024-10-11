package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/YaroslavGaponov/grotto/internal/chunkservice"
	"github.com/YaroslavGaponov/grotto/internal/configure"
	"github.com/YaroslavGaponov/grotto/pkg/logger"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	config := configure.NewConfigure()
	log := logger.NewLogger(config.LogLevel)

	log.Infof("chunk service is starting...")
	chunk := chunkservice.New(config)
	go chunk.Start()

	<-sigs

	chunk.Stop()
}
