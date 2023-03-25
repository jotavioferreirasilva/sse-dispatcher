package handler

import (
	"backend/src/sse"
	sse_server "backend/src/sse_server/proto"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PushMessageService struct {
	sse_server.PushMessageServer
}

func (pushMessageService PushMessageService) PushMessage(ctx context.Context, req *sse_server.PushMessageRequest) (*empty.Empty, error) {
	if req.Message == "" {
		err := status.New(codes.InvalidArgument, "Message cannot be empty").Err()
		return nil, err
	}

	sse.Clients.Publish(req.Message)

	return &empty.Empty{}, nil
}
