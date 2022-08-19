package limiter

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTimeout          = errors.New("context timed out")
	ErrBurstExceeded    = errors.New("burst exceeded limit")
	ErrInvalidBurstSize = errors.New("invalid burst size")
)

type Limiter interface {
	Reserve(ctx context.Context) (*Reservation, error)
	ReserveN(ctx context.Context, n int) (*Reservation, error)
	Wait(ctx context.Context) error
	WaitN(ctx context.Context, n int) error
}

type Limit struct {
	Count    int
	Interval time.Duration
}
