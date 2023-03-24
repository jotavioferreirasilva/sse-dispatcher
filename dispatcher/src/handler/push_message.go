package handler

import (
	"context"
	dispatcher_server "dispatcher/src/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type PushMessageService struct {
	dispatcher_server.PushMessageServer
}

func (pushMessageService PushMessageService) PushMessage(ctx context.Context, req *dispatcher_server.PushMessageRequest) (*empty.Empty, error) {
	if req.Message == "" {
		err := status.New(codes.InvalidArgument, "Message cannot be empty").Err()
		return nil, err
	}

	log.Printf("Message: %s", req.Message)

	return &empty.Empty{}, nil
}
