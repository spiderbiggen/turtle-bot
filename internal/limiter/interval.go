package limiter

import (
	"context"
	"sync"
	"time"
)

type IntervalWindow struct {
	Limits []Limit
	count  map[time.Duration]uint
	offset map[time.Duration]uint
	start  time.Time
	mu     *sync.Mutex
}

func NewIntervalWindow(limits ...Limit) *IntervalWindow {
	return &IntervalWindow{
		Limits: limits,
		count:  make(map[time.Duration]uint),
		offset: make(map[time.Duration]uint),
		mu:     &sync.Mutex{},
	}
}

func (i *IntervalWindow) Wait(ctx context.Context) (Incrementer, error) { return i.WaitBurst(ctx, 1) }

func (i *IntervalWindow) WaitBurst(ctx context.Context, n uint) (Incrementer, error) {
	now := time.Now()
	c := make(chan error)
	go func() {
		i.mu.Lock()

		i.resetStart(now)
		i.resetIntervals(now)
		max, err := i.maxIntervalN(n)
		if err != nil || max == 0 {
			c <- err
			return
		}
		offset, _ := i.offset[max]
		t := time.NewTimer(i.start.Add(time.Duration(offset+1) * max).Sub(now))
		defer t.Stop()
		select {
		case <-t.C:
			c <- nil
		case <-ctx.Done():
		}
	}()

	select {
	case err := <-c:
		if err != nil {
			i.mu.Unlock()
			return func() {}, nil
		}
		return func() {
			i.resetIntervals(time.Now())
			i.incrementN(n)
			i.mu.Unlock()
		}, nil
	case <-ctx.Done():
		i.mu.Unlock()
		return func() {}, ErrTimeout
	}
}

func (i *IntervalWindow) maxIntervalN(n uint) (time.Duration, error) {
	var max time.Duration
	for _, l := range i.Limits {
		if n > l.Count {
			return 0, ErrBurstExceeded
		}
		d, _ := i.count[l.Interval]
		if d+n > l.Count {
			if l.Interval > max {
				max = l.Interval
			}
		}
	}
	return max, nil
}

func (i *IntervalWindow) resetStart(t time.Time) {
	allExpired := true
	for _, l := range i.Limits {
		offset, _ := i.offset[l.Interval]
		duration := time.Duration(offset+1) * l.Interval
		expires := i.start.Add(duration)
		allExpired = allExpired && expires.Before(t)
	}
	if allExpired {
		i.start = time.Now().Truncate(time.Second).Add(time.Second)
		for _, l := range i.Limits {
			i.offset[l.Interval] = 0
			i.count[l.Interval] = 0
		}
	}
}

func (i *IntervalWindow) resetIntervals(t time.Time) {
	for _, l := range i.Limits {
		offset, _ := i.offset[l.Interval]
		for i.start.Add(time.Duration(offset+1) * l.Interval).Before(t) {
			i.offset[l.Interval] = offset + 1
			i.count[l.Interval] = 0
			offset++
		}
	}
}

func (i *IntervalWindow) incrementN(n uint) {
	for _, l := range i.Limits {
		i.count[l.Interval] += n
	}
}
