package configure

import (
	"os"
	"strconv"
	"strings"
)

type Configure struct {
	LogLevel          string
	ChunkServiceAddr  string
	MasterServiceAddr string
	ChunkServiceUrls  []string
	DriverType        string
	DiskDriverRootDir string
	ChunkReplicas     int
}

func NewConfigure() Configure {
	return Configure{
		LogLevel:          getEnvOrDefault("LOG_LEVEL", "all"),
		ChunkServiceAddr:  getEnvOrDefault("CHUNK_SERVICE_ADDR", ":8001"),
		MasterServiceAddr: getEnvOrDefault("MASTER_SERVICE_ADDR", ":8000"),
		ChunkServiceUrls:  strings.Split(getEnvOrDefault("CHUNK_SERVICE_URLS", ""), ","),
		DriverType:        getEnvOrDefault("DRIVER_TYPE", "memory"), // memory or disk
		DiskDriverRootDir: getEnvOrDefault("DRIVER_DISK_ROOT", "/var/grotto"),
		ChunkReplicas:     getEnvOrDefaultInt("CHUNK_REPLICAS", 3),
	}
}

func getEnvOrDefault(name, defValue string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defValue
	}
	return value
}

func getEnvOrDefaultInt(name string, defValue int) int {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defValue
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return defValue
	}
	return num
}
