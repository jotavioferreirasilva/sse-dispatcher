package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var DispatcherServerConfiguration Config

type Config struct {
	Host                 string
	Port                 uint32
	DockerNetwork        string
	DockerSSEServerLabel string
	DockerSSEServerPort  uint32
}

func NewConfig(path string) Config {
	if strings.TrimSpace(path) != "" {
		loadEnvironmentFile(path)
	}

	var config Config
	config.Host = getDispatcherHost()
	config.Port = getDispatcherPort()
	config.DockerNetwork = getDockerNetwork()
	config.DockerSSEServerLabel = getDockerSSEServerLabel()
	config.DockerSSEServerPort = getDockerSSEServerPort()
	return config
}

func getDispatcherHost() string {
	return os.Getenv("HOST")
}

func getDispatcherPort() uint32 {
	envDispatcherPort := os.Getenv("PORT")

	dispatcherPort, err := strconv.Atoi(envDispatcherPort)

	if err != nil {
		log.Fatal(err.Error())
	}

	return uint32(dispatcherPort)
}

func getDockerNetwork() string {
	return os.Getenv("DOCKER_NETWORK")
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

func getDockerSSEServerLabel() string {
	return os.Getenv("DOCKER_SSE_SERVER_LABEL")
}

func getDockerSSEServerPort() uint32 {
	envDockerSSEServerPort := os.Getenv("DOCKER_SSE_SERVER_PORT")

	dockerSSEServerPort, err := strconv.Atoi(envDockerSSEServerPort)

	if err != nil {
		log.Fatal(err.Error())
	}

	return uint32(dockerSSEServerPort)
}
