package limiter

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTimeout       = errors.New("context timed out")
	ErrBurstExceeded = errors.New("burst exceeded limit")
)

type Limiter interface {
	Wait(ctx context.Context) (Incrementer, error)
	WaitBurst(ctx context.Context, n uint) (Incrementer, error)
}

type Limit struct {
	Count    uint
	Interval time.Duration
}

type Incrementer func()
