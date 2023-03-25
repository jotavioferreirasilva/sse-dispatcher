package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var BackendConfiguration Config

type Config struct {
	DispatcherServerHost string
	DispatcherServerPort uint32
	SSEServerHost        string
	SSEServerPort        uint32
}

func NewConfig(path string) Config {
	if strings.TrimSpace(path) != "" {
		loadEnvironmentFile(path)
	}

	var config Config
	config.DispatcherServerHost = getDispatcherHost()
	config.DispatcherServerPort = getDispatcherPort()
	config.SSEServerHost = getSSEServerHost()
	config.SSEServerPort = getSSEServerPort()
	return config
}

func getDispatcherPort() uint32 {
	envDispatcherPort := os.Getenv("DISPATCHER_PORT")

	dispatcherPort, err := strconv.Atoi(envDispatcherPort)

	if err != nil {
		log.Fatal(err.Error())
	}

	return uint32(dispatcherPort)
}

func getDispatcherHost() string {
	return os.Getenv("DISPATCHER_HOST")
}

func loadEnvironmentFile(path string) {
	log.Println("Trying to load file with environment variables")

	_, err := os.Stat(path)
	if err != nil {
		log.Println("File not found. Values will be loaded from machine")
		return
	}

	log.Println("File found, loading values")
	if err := godotenv.Load(path); err != nil {
		log.Fatal(err.Error())
	}
}

func getSSEServerPort() uint32 {
	envSSEServerPort := os.Getenv("SSE_SERVER_PORT")

	sseServerPort, err := strconv.Atoi(envSSEServerPort)

	if err != nil {
		log.Fatal(err.Error())
	}

	return uint32(sseServerPort)
}

func getSSEServerHost() string {
	return os.Getenv("SSE_SERVER_HOST")
}
