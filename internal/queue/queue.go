package queue

import (
	"errors"
	"sync"
)

var (
	ErrEmpty            = errors.New("queue was empty")
	ErrCapacityExceeded = errors.New("queue capacity exceeds limits")
	ErrIndexOutOfRange  = errors.New("index out of range")
)

// Queue is a concurrency-safe rotating queue implementation.
type Queue[T any] struct {
	Data     []*T
	capacity int
	current  int
	size     int
	mu       *sync.RWMutex
}

func New[T any](cap int) Queue[T] {
	return Queue[T]{
		Data:     make([]*T, cap),
		capacity: cap,
		mu:       &sync.RWMutex{},
	}
}

func (q *Queue[T]) Len() int               { return q.size }
func (q *Queue[T]) Cap() int               { return q.capacity }
func (q *Queue[T]) rollOver(index int) int { return index % q.capacity }
func (q *Queue[T]) index(offset int) int   { return q.rollOver(q.current + offset) }

func (q *Queue[T]) Push(data T) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.size >= q.capacity {
		return ErrCapacityExceeded
	}
	i := q.index(q.size)
	q.Data[i] = &data
	q.size++
	return nil
}

func (q *Queue[T]) Pop() (*T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	el := q.Data[q.current]
	if el == nil {
		return nil, ErrEmpty
	}
	q.Data[q.current] = nil
	q.size--
	if q.size == 0 {
		q.current = 0
	} else {
		q.current = q.index(1)
	}
	return el, nil
}

func (q *Queue[T]) Peek() (*T, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	el := q.Data[q.current]
	if el == nil {
		return nil, ErrEmpty
	}
	return el, nil
}

func (q *Queue[T]) TryPeek() *T {
	a, _ := q.Peek()
	return a
}

func (q *Queue[T]) PeekIndex(index int) (*T, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if index < 0 || index >= q.size {
		return nil, ErrIndexOutOfRange
	}
	i := q.rollOver(q.current + index)
	return q.Data[i], nil
}

func (q *Queue[T]) Remaining() []*T {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if q.size == 0 {
		return nil
	}
	result := make([]*T, q.size)
	m, r := q.current+q.size, 0
	if m > q.capacity {
		r = m - q.capacity
		m = q.capacity
	}

	copy(result, q.Data[q.current:m])
	if r > 0 {
		copy(result[m-q.current:], q.Data[:r])
	}
	return result
}
