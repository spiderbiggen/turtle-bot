package queue

import (
	"errors"
	"sync"
)

var (
	ErrNoElements = errors.New("no elements in queue")
)

// Queue is a concurrency-safe rotating queue implementation.
type Queue[T any] struct {
	data     []*T
	baseSize int
	mu       *sync.RWMutex
}

func New[T any](base int) Queue[T] {
	return Queue[T]{
		data:     make([]*T, 0, base),
		baseSize: base,
		mu:       &sync.RWMutex{},
	}
}

func (q *Queue[T]) Len() int { return len(q.data) }

func (q *Queue[T]) Append(data *T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, data)
}

func (q *Queue[T]) AppendMany(data []*T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, data...)
}

func (q *Queue[T]) Pop() (*T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		return nil, ErrNoElements
	}
	var result *T
	result, q.data = q.data[0], q.data[1:]
	if len(q.data) == 0 {
		q.data = make([]*T, 0, q.baseSize)
	}
	return result, nil
}

// PopN pops up to n elements from the queue
func (q *Queue[T]) PopN(n int) []*T {
	q.mu.Lock()
	defer q.mu.Unlock()
	l := len(q.data)
	if l == 0 {
		return nil
	}
	if n > l {
		n = l
	}
	var result []*T
	result, q.data = q.data[:n], q.data[n:]
	if len(q.data) == 0 {
		q.data = make([]*T, 0, q.baseSize)
	}
	return result
}
