package broker

import (
	"context"
	"time"

	"github.com/ruslan-onishchenko/go-test-task/pkg/queue"
)

type Queue interface {
	Push(value any) bool
	Get(ctx context.Context, timeout time.Duration) (any, bool)
	Run(ctx context.Context)
}

type broker struct {
	initQueue func() Queue
	queueMap  map[string]Queue
	maxQueue  int
}

func (b *broker) Push(queueName string, value any) bool {
	if _, ok := b.queueMap[queueName]; !ok {
		if b.maxQueue != 0 && len(b.queueMap) >= b.maxQueue {
			return false
		}

		b.queueMap[queueName] = b.initQueue()
	}

	return b.queueMap[queueName].Push(value)
}

func (b *broker) Get(ctx context.Context, queueName string, timeout time.Duration) (any, bool) {
	if _, ok := b.queueMap[queueName]; !ok {
		return nil, false
	}

	return b.queueMap[queueName].Get(ctx, timeout)
}

type Broker interface {
	Push(queueName string, value any) bool
	Get(ctx context.Context, queueName string, timeout time.Duration) (any, bool)
}

type Option func(*broker)

func QueueEmpty() Option {
	return func(b *broker) {
		b.initQueue = func() Queue { return queue.NewEmpty() }
	}
}

func QueueSimple() Option {
	return func(b *broker) {
		b.initQueue = func() Queue { return queue.NewSimple() }
	}
}

func QueueFifo(ctx context.Context, opts ...queue.FifoOption) Option {
	return func(b *broker) {
		b.initQueue = func() Queue { return queue.NewFifo(ctx, opts...) }
	}
}

func CustomQueue(fun func() Queue) Option {
	return func(b *broker) {
		b.initQueue = fun
	}
}

func MaxQueue(maxQueue int) Option {
	return func(b *broker) {
		b.maxQueue = maxQueue
	}
}

func New(opts ...Option) Broker {
	b := broker{
		initQueue: func() Queue { return queue.NewEmpty() },
		queueMap:  map[string]Queue{},
		maxQueue:  0,
	}

	for _, opt := range opts {
		opt(&b)
	}

	return &b
}
