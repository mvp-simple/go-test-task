package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ruslan-onishchenko/go-test-task/pkg/broker"
	"github.com/ruslan-onishchenko/go-test-task/pkg/queue"
	"github.com/ruslan-onishchenko/go-test-task/pkg/servelayer"
)

var (
	maxCountQueueMsg    = flag.Int("max-queue-message", 0, "maximum count of message at one queue")
	maxCountQueue       = flag.Int("max-queue-count", 0, "maximum count of queue")
	opGetTimeout        = flag.Int("operation-get-timeout", 0, "timeout of get from queue")
	opGetDisableTimeout = flag.Bool("operation-get-disable-timeout", false, "disable timeout of get from queue")
	port                = flag.Uint("port", 8081, "disable timeout of get from queue")
	grpcPort            = flag.Uint("grpc-port", 12201, "disable timeout of get from queue")
)

func main() {
	flag.Parse()

	if *port == *grpcPort {
		log.Fatal("port can not be equal for grpc port")
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		select {
		case <-exit:
			cancel()
		case <-ctx.Done():
		}
	}()

	run(ctx)
}

func run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		queueOpts  []queue.FifoOption
		brokerOpts []broker.Option
	)

	if *maxCountQueueMsg > 0 {
		queueOpts = append(queueOpts, queue.FifoMaxLen(*maxCountQueueMsg))
	}

	if *opGetTimeout != 0 && !(*opGetDisableTimeout) {
		queueOpts = append(queueOpts, queue.FifoTimeOut(time.Duration(*opGetTimeout)*time.Second))
	}

	if *opGetDisableTimeout {
		queueOpts = append(queueOpts, queue.FifoDisableTimeout())
	}

	if *maxCountQueue > 0 {
		brokerOpts = append(brokerOpts, broker.MaxQueue(*maxCountQueue))
	}

	brokerOpts = append(brokerOpts, broker.QueueFifo(ctx, queueOpts...))

	queueServiceServer, err := servelayer.New(broker.New(brokerOpts...))
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var errGateway, errGrpc error
	go func() {
		defer wg.Done()
		if errGateway = runGateway(ctx, *grpcPort, *port); errGateway != nil {
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if errGrpc = runGrpc(ctx, queueServiceServer, *grpcPort); errGrpc != nil {
			cancel()
		}
	}()
	wg.Wait()

	if err := errors.Join(errGateway, errGrpc); err != nil {
		log.Fatal(err)
	}
}
