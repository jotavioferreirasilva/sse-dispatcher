package server

import (
	"dispatcher/src/config"
	"dispatcher/src/handler"
	dispatcher_server "dispatcher/src/proto"
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
	dispatcher_server.RegisterPushMessageServer(server.Server, &handler.PushMessageService{})
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
	addr := fmt.Sprintf("%s:%d", config.DispatcherServerConfiguration.Host, config.DispatcherServerConfiguration.Port)

	s.Listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %s. Error: %s", addr, err)
	}

	log.Printf("gRPC Server: Started listening at %s", addr)
	return
}
