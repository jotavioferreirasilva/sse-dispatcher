package handler

import (
	"context"
	"dispatcher/src/config"
	"dispatcher/src/docker_client"
	dispatcher_server "dispatcher/src/proto"
	sse_server "dispatcher/src/sse_client/sse_server_proto"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

type PushMessageService struct {
	dispatcher_server.PushMessageServer
}

func (pushMessageService PushMessageService) PushMessage(context context.Context, req *dispatcher_server.PushMessageRequest) (*empty.Empty, error) {
	if req.Message == "" {
		err := status.New(codes.InvalidArgument, "Message cannot be empty").Err()
		return nil, err
	}

	containerFilter := filters.NewArgs()
	containerFilter.Add("label", config.DispatcherServerConfiguration.DockerSSEServerLabel)

	containers, err := docker_client.DockerClient.ContainerList(context, types.ContainerListOptions{Filters: containerFilter})

	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		if err != nil {
			return nil, err
		}
		ipAddress := container.NetworkSettings.Networks[config.DispatcherServerConfiguration.DockerNetwork].IPAddress

		if ipAddress == "" {
			log.Printf("container ip address is empty: %s", container.Names)
			return nil, err
		}

		err := sendMessageToBackend(context, req, ipAddress)

		if err != nil {
			log.Print(err)
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

func sendMessageToBackend(context context.Context, req *dispatcher_server.PushMessageRequest, ipAddress string) error {
	host := fmt.Sprintf("%s:%d", ipAddress, config.DispatcherServerConfiguration.DockerSSEServerPort)
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Print(err)
		}
	}(conn)
	c := sse_server.NewPushMessageClient(conn)

	log.Printf(">>> Sending message \"%s\" to %s...>>>\n", req.Message, host)
	_, err = c.PushMessage(context, &sse_server.PushMessageRequest{
		Message: req.Message,
	})

	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
