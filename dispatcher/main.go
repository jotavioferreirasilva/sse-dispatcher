package main

import (
	"dispatcher/src/config"
	"dispatcher/src/server"
)

func init() {
	config.DispatcherServerConfiguration = config.NewConfig(".env")
	server.GRPCServer = server.Start()
}

func main() {
	server.GRPCServer.Server.Serve(server.GRPCServer.Listener)
}
