package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/ruslan-onishchenko/go-test-task/pkg/servelayer/queue/v1"
)

func runGrpc(ctx context.Context, server pb.QueueServiceServer, grpcPort uint) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcPort))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterQueueServiceServer(s, server)

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return s.Serve(lis)
}
