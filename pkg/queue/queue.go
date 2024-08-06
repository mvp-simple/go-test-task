package queue

import (
	"context"
	"time"
)

type Queue interface {
	Push(value any) bool
	Get(ctx context.Context, timeout time.Duration) (any, bool)
	Run(ctx context.Context)
}
