package timeouts

import (
	"context"
	"time"
)

type Executor struct {
	timeout time.Duration
}

func (e *Executor) Execute(fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	return fn(ctx)
}
