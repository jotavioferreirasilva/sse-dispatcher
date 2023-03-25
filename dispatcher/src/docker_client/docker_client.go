package docker_client

import (
	"github.com/docker/docker/client"
	"log"
)

var DockerClient *client.Client

func NewDockerClient() *client.Client {
	dockerClient, err := client.NewClientWithOpts(client.WithVersion("1.41"))

	if err != nil {
		log.Fatalf("Failed to start docker client: %s", err)
	}

	return dockerClient
}
