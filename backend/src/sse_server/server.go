package server

import (
	"backend/src/config"
	"backend/src/sse_server/handler"
	sse_server "backend/src/sse_server/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var GRPCServer *GRPC

type GRPC struct {
	Server   *grpc.Server
	Listener net.Listener
}

func Start() *GRPC {
	server := &GRPC{
		Server: grpc.NewServer(),
	}
	server.initListener()
	sse_server.RegisterPushMessageServer(server.Server, &handler.PushMessageService{})
	go stopListener(server.Server)

	return server
}

func stopListener(server *grpc.Server) {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	_ = <-stopSignal

	log.Println("Stopping gRPC server...")
	server.GracefulStop()
}

func (s *GRPC) initListener() {
	var err error
	addr := fmt.Sprintf("%s:%d", config.BackendConfiguration.SSEServerHost, config.BackendConfiguration.SSEServerPort)

	s.Listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %s. Error: %s", addr, err)
	}

	log.Printf("Started listening at %s", addr)
	return
}
