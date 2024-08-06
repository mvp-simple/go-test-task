package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ruslan-onishchenko/go-test-task/pkg/servelayer/queue/v1"
)

func runGateway(ctx context.Context, grpcPort, port uint) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := pb.RegisterQueueServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%v", grpcPort),
		opts,
	)
	if err != nil {
		return err
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux); err != nil {
		return err
	}

	return nil
}
