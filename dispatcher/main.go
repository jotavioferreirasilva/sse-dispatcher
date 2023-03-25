package main

import (
	"dispatcher/src/config"
	"dispatcher/src/docker_client"
	"dispatcher/src/server"
)

func init() {
	config.DispatcherServerConfiguration = config.NewConfig(".env")
	server.GRPCServer = server.Start()
	docker_client.DockerClient = docker_client.NewDockerClient()
}

func main() {
	server.GRPCServer.Server.Serve(server.GRPCServer.Listener)
}
