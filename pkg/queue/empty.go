package queue

import (
	"context"
	"time"
)

type empty struct{}

func (q *empty) Push(any) bool {
	return false
}

func (q *empty) Get(context.Context, time.Duration) (any, bool) {
	return nil, false
}

func (q *empty) Run(context.Context) {}

func NewEmpty() Queue {
	return &empty{}
}
