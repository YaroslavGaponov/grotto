package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/YaroslavGaponov/grotto/internal/chunkservice"
	"github.com/YaroslavGaponov/grotto/internal/configure"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	config := configure.NewConfigure()

	fmt.Println("chunk service is starting...")
	chunk := chunkservice.New(config)
	go chunk.Start()

	<-sigs

	chunk.Stop()
}
