package limiter

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	"weeb_bot/internal/queue"
)

type RollingWindow struct {
	Limits []Limit
	count  map[time.Duration]queue.Queue[time.Time]
	mu     *sync.Mutex
}

func NewRollingWindow(limits ...Limit) *RollingWindow {
	m := make(map[time.Duration]queue.Queue[time.Time])
	for _, limit := range limits {
		m[limit.Interval] = queue.New[time.Time](int(limit.Count))
	}
	return &RollingWindow{
		Limits: limits,
		count:  m,
		mu:     &sync.Mutex{},
	}
}

func (i *RollingWindow) Wait(ctx context.Context) (Incrementer, error) { return i.WaitBurst(ctx, 1) }

func (i *RollingWindow) WaitBurst(ctx context.Context, n uint) (Incrementer, error) {
	now := time.Now()
	c := make(chan error)

	go func() {
		i.mu.Lock()

		if ctx.Err() != nil {
			return
		}
		_ = i.resetIntervals(now)
		max, err := i.maxIntervalN(n)
		if err != nil || max == 0 {
			c <- err
			return
		}
		t := time.NewTimer(max)
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
			return func() {}, err
		}
		return func() {
			_ = i.resetIntervals(time.Now())
			i.incrementN(n)
			i.mu.Unlock()
		}, nil
	case <-ctx.Done():
		i.mu.Unlock()
		return func() {}, ctx.Err()
	}
}

func (i *RollingWindow) maxIntervalN(n uint) (time.Duration, error) {
	var max time.Duration
	for _, l := range i.Limits {
		if n > l.Count {
			return 0, ErrBurstExceeded
		}
		d, ok := i.count[l.Interval]
		if !ok {
			d = queue.New[time.Time](int(l.Count))
			i.count[l.Interval] = d
		}

		if uint(d.Len())+n > l.Count {
			if l.Interval > max {
				max = l.Interval
			}
		}
	}
	return max, nil
}

func (i *RollingWindow) resetIntervals(t time.Time) error {
	for _, l := range i.Limits {
		q := i.count[l.Interval]
		if q.Len() == 0 {
			continue
		}

		expires := t.Add(-l.Interval)
		w, err := q.Peek()
		if err != nil {
			return err
		}
		for w != nil && w.Before(expires) {
			_, _ = q.Pop()
			w, err = q.Peek()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *RollingWindow) incrementN(n uint) {
	t := time.Now()
	for _, limit := range i.Limits {
		q := i.count[limit.Interval]
		for i := uint(0); i < n; i++ {
			err := q.Push(t)
			if err != nil {
				log.Error(err)
				break
			}
		}
	}
}
