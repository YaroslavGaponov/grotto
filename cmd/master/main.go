package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/YaroslavGaponov/grotto/internal/configure"
	"github.com/YaroslavGaponov/grotto/internal/masterservice"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	config := configure.NewConfigure()

	fmt.Println("master service is starting...")
	master := masterservice.New(config.MasterServiceAddr, config.ChunkServiceUrls)
	go master.Start()

	<-sigs

	master.Stop()
}
