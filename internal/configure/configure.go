package configure

import (
	"os"
	"strings"
)

type Configure struct {
	ChunkServiceAddr  string
	MasterServiceAddr string
	ChunkServiceUrls  []string
}

func NewConfigure() Configure {
	return Configure{
		ChunkServiceAddr:  getEnvOrDefault("CHUNK_SERVICE_ADDR", ":8001"),
		MasterServiceAddr: getEnvOrDefault("MASTER_SERVICE_ADDR", ":8000"),
		ChunkServiceUrls:  strings.Split(getEnvOrDefault("CHUNK_SERVICE_URLS", ""), ","),
	}
}

func getEnvOrDefault(name, defValue string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defValue
	}
	return value
}
