package queue

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type fifo struct {
	mu           *sync.Mutex
	valueList    *list.List
	listenerList *list.List

	maxLen         int
	timeout        time.Duration
	disableTimeout bool
}

func (q *fifo) Push(value any) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.maxLen != 0 && q.valueList.Len() >= q.maxLen {
		return false
	}

	q.valueList.PushBack(value)
	return true
}

func (q *fifo) Get(ctx context.Context, timeout time.Duration) (any, bool) {
	var (
		result any
		ok     bool

		timeOutCtx       context.Context
		cancelTimeOutCtx context.CancelFunc
	)

	if !q.disableTimeout && (q.timeout != 0 || timeout != 0) {
		if timeout > q.timeout || timeout == 0 {
			timeout = q.timeout
		}

		timeOutCtx, cancelTimeOutCtx = context.WithTimeout(ctx, timeout)
	} else {
		timeOutCtx, cancelTimeOutCtx = context.WithCancel(ctx)
	}
	defer cancelTimeOutCtx()

	listenerCtx, cancelListenCtx := context.WithCancel(context.TODO())
	defer cancelListenCtx()

	q.mu.Lock()
	listener := q.listenerList.PushBack(func(data any) {
		result = data
		ok = true
		cancelListenCtx()
	})
	q.mu.Unlock()

	select {
	case <-timeOutCtx.Done():
		q.mu.Lock()
		q.listenerList.Remove(listener)
		q.mu.Unlock()
	case <-listenerCtx.Done():
	}

	return result, ok
}

func (q *fifo) Run(ctx context.Context) {
	go q.run(ctx)
}

func (q *fifo) run(ctx context.Context) {
	ticker := time.NewTicker(time.Nanosecond)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q.tick()
		}
	}
}

func (q *fifo) tick() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.valueList.Len() == 0 || q.listenerList.Len() == 0 {
		return
	}

	valueElement := q.valueList.Back()
	value := q.valueList.Remove(valueElement)

	listenerElement := q.listenerList.Back()
	listener := q.listenerList.Remove(listenerElement)

	typeFunc := listener.(func(any))
	typeFunc(value)
}

type FifoOption func(*fifo)

func FifoTimeOut(timeOut time.Duration) FifoOption {
	return func(q *fifo) {
		q.timeout = timeOut
	}
}

func FifoMaxLen(maxLen int) FifoOption {
	return func(q *fifo) {
		q.maxLen = maxLen
	}
}

func FifoDisableTimeout() FifoOption {
	return func(q *fifo) {
		q.disableTimeout = true
	}
}

func NewFifo(ctx context.Context, opts ...FifoOption) Queue {
	q := fifo{
		valueList:    list.New(),
		listenerList: list.New(),
		mu:           &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(&q)
	}

	q.Run(ctx)
	return &q
}
