package limiter

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
	"time"
)

type IntervalWindow struct {
	Limits       []Limit
	maxInterval  time.Duration
	maxBurst     int
	reservations []*Reservation
	counts       map[time.Duration]int
	windowStart  map[time.Duration]time.Time
	mu           *sync.Mutex
	cleanup      bool
}

func NewIntervalWindow(limits ...Limit) *IntervalWindow {
	var maxInterval time.Duration
	var minBurst = math.MaxInt
	for _, l := range limits {
		if l.Interval > maxInterval {
			maxInterval = l.Interval
		}
		if l.Count < minBurst {
			minBurst = l.Count
		}
	}

	return &IntervalWindow{
		Limits:      limits,
		maxInterval: maxInterval,
		maxBurst:    minBurst,
		counts:      make(map[time.Duration]int),
		windowStart: make(map[time.Duration]time.Time),
		mu:          &sync.Mutex{},
	}
}

func (i *IntervalWindow) StartCleanup() {
	i.mu.Lock()
	defer i.mu.Unlock()
	if !i.cleanup {
		go func() {
			t := time.NewTicker(5 * time.Minute)
			for {
				select {
				case <-t.C:
					count := i.Cleanup()
					log.Debugf("Cleaned up %d old reservations", count)
				}
			}
		}()
	}
}

func (i *IntervalWindow) Reserve(ctx context.Context) (*Reservation, error) {
	return i.ReserveN(ctx, 1)
}

func (i *IntervalWindow) ReserveN(ctx context.Context, n int) (*Reservation, error) {
	if n <= 0 {
		return nil, ErrInvalidBurstSize
	}
	if n > i.maxBurst {
		return nil, ErrBurstExceeded
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	now := time.Now()
	var reservation *Reservation
	if len(i.reservations) == 0 {
		reservation = &Reservation{count: n, time: now, ctx: ctx}
		i.reservations = append(i.reservations, reservation)
		for _, l := range i.Limits {
			i.counts[l.Interval] = n
			i.windowStart[l.Interval] = now
		}
		return reservation, nil
	} else {
		var earliestAvailableStart = now
		for _, l := range i.Limits {
			s := i.windowStart[l.Interval]
			if c, _ := i.counts[l.Interval]; c+n > l.Count {
				s = s.Add(l.Interval)
			}
			if s.After(earliestAvailableStart) {
				earliestAvailableStart = s
			}
		}
		if dl, ok := ctx.Deadline(); ok && earliestAvailableStart.After(dl) {
			return nil, ErrTimeout
		}
		reservation = &Reservation{count: n, time: earliestAvailableStart, ctx: ctx}
		i.reservations = append(i.reservations, reservation)
		for _, l := range i.Limits {
			if earliestAvailableStart.After(i.windowStart[l.Interval].Add(l.Interval)) {
				i.counts[l.Interval] = n
				i.windowStart[l.Interval] = earliestAvailableStart
			} else {
				i.counts[l.Interval] += n
			}
		}
	}
	return reservation, nil
}

func (i *IntervalWindow) Wait(ctx context.Context) error { return i.WaitN(ctx, 1) }

func (i *IntervalWindow) WaitN(ctx context.Context, n int) error {
	r, err := i.ReserveN(ctx, n)
	if err != nil {
		return err
	}
	timer := time.NewTimer(r.Delay())
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-r.ctx.Done():
		return ErrTimeout
	}
}

func (i *IntervalWindow) Cleanup() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	if len(i.reservations) == 0 {
		return 0
	}
	now := time.Now()
	temp := make([]*Reservation, 0, len(i.reservations)/2)
	for _, reservation := range i.reservations {
		for _, l := range i.Limits {
			if now.Before(reservation.time.Add(2 * l.Interval)) {
				temp = append(temp, reservation)
				break
			}
		}
	}
	count := len(i.reservations) - len(temp)
	i.reservations = temp
	return count
}
