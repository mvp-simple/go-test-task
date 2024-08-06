package servelayer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ruslan-onishchenko/go-test-task/pkg/broker"
	pb "github.com/ruslan-onishchenko/go-test-task/pkg/servelayer/queue/v1"
)

var _ pb.QueueServiceServer = &queueService{}

func New(b broker.Broker) (pb.QueueServiceServer, error) {
	if b == nil {
		return nil, errors.New("broker.Broker is empty value")
	}

	return &queueService{
		broker: b,
	}, nil
}

type queueService struct {
	pb.UnimplementedQueueServiceServer

	broker broker.Broker
}

func (q *queueService) Push(ctx context.Context, request *pb.PushRequest) (*pb.PushResponse, error) {
	if request.GetMessage() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing parameter message")
	}

	if !q.broker.Push(request.GetQueue(), request.GetMessage()) {
		return nil, status.Error(codes.InvalidArgument, "can not push new message")
	}

	return &pb.PushResponse{}, nil
}

func (q *queueService) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	val, ok := q.broker.Get(ctx, request.GetQueue(), time.Duration(request.GetTimeout())*time.Second)
	if !ok {
		return nil, status.Error(codes.NotFound, "404")
	}

	return &pb.GetResponse{
		Message: fmt.Sprint(val),
	}, nil
}
