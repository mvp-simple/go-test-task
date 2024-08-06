package queue

import (
	"container/list"
	"context"
	"time"
)

type simple struct {
	list *list.List
}

func (q *simple) Push(value any) bool {
	q.list.PushBack(value)
	return true
}

func (q *simple) Get(context.Context, time.Duration) (any, bool) {
	ok := q.list.Len() != 0
	back := q.list.Back()
	return back, ok
}

func (q *simple) Run(ctx context.Context) {}

func NewSimple() Queue {
	return &simple{list: list.New()}
}
