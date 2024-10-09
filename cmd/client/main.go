package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/YaroslavGaponov/grotto/pkg/client"
	"github.com/YaroslavGaponov/grotto/pkg/common"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		help()
		return
	}
	switch args[0] {
	case "upload":
		if len(args) != 3 {
			help()
			return
		}
		upload(args[1], args[2])
	case "download":
		if len(args) != 3 {
			help()
			return
		}
		download(args[1], args[2])
	case "catalog":
		if len(args) != 2 {
			help()
			return
		}
		catalog(args[1])
	case "watch":
		if len(args) != 2 {
			help()
			return
		}
		watch(args[1])
	case "help":
	default:
		help()
	}

}

func help() {
	fmt.Println("Usage:")
	fmt.Println("\tclient [command] args...")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("\thelp                   Print help")
	fmt.Println("\tupload {file} {url}    Upload local file to grotto")
	fmt.Println("\tdownload {file} {url}  Download file from grotto")
	fmt.Println("\tcatalog {url}          Print catalog")
	fmt.Println("\twatch {url}            Print all events")
}

func upload(file, url string) {
	fmt.Println("uploading ", file)
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	name := filepath.Base(file)
	c := client.NewClient(url)
	err = c.Save(name, data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("done")
}

func download(file, url string) {
	fmt.Println("downloading ", file)
	name := filepath.Base(file)
	c := client.NewClient(url)
	data, err := c.Load(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("done")
}

func catalog(url string) {
	c := client.NewClient(url)
	list, err := c.List()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, name := range list {
		fmt.Println(name)
	}
}

func watch(url string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	c := client.NewClient(url)
	channel := make(chan common.Event)

	fmt.Println("Waiting for events...")
	go c.Watch(channel)
	for {
		select {
		case <-sigs:
			fmt.Println("bye")
			close(channel)
			return
		case event := <-channel:
			fmt.Println(event.Action, event.File)
		}
	}
}
