package limiter

import (
	"context"
	"time"
)

type Reservation struct {
	count int
	time  time.Time
	ctx   context.Context
}

func (r *Reservation) Delay() time.Duration {
	return r.time.Sub(time.Now())
}
