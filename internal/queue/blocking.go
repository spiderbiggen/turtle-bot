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

// BlockingQueue is a concurrency-safe rotating queue implementation.
type BlockingQueue[T any] struct {
	data     []*T
	capacity int
	current  int
	size     int
	mu       *sync.RWMutex
}

func NewBlocking[T any](cap int) BlockingQueue[T] {
	return BlockingQueue[T]{
		data:     make([]*T, cap),
		capacity: cap,
		mu:       &sync.RWMutex{},
	}
}

func (q *BlockingQueue[T]) Len() int               { return q.size }
func (q *BlockingQueue[T]) Cap() int               { return q.capacity }
func (q *BlockingQueue[T]) rollOver(index int) int { return index % q.capacity }
func (q *BlockingQueue[T]) index(offset int) int   { return q.rollOver(q.current + offset) }

func (q *BlockingQueue[T]) Push(data T) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.size >= q.capacity {
		return ErrCapacityExceeded
	}
	i := q.index(q.size)
	q.data[i] = &data
	q.size++
	return nil
}

func (q *BlockingQueue[T]) Pop() (*T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	el := q.data[q.current]
	if el == nil {
		return nil, ErrEmpty
	}
	q.data[q.current] = nil
	q.size--
	if q.size == 0 {
		q.current = 0
	} else {
		q.current = q.index(1)
	}
	return el, nil
}

func (q *BlockingQueue[T]) Peek() (*T, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	el := q.data[q.current]
	if el == nil {
		return nil, ErrEmpty
	}
	return el, nil
}

func (q *BlockingQueue[T]) TryPeek() *T {
	a, _ := q.Peek()
	return a
}

func (q *BlockingQueue[T]) PeekIndex(index int) (*T, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if index < 0 || index >= q.size {
		return nil, ErrIndexOutOfRange
	}
	i := q.rollOver(q.current + index)
	return q.data[i], nil
}

func (q *BlockingQueue[T]) Remaining() []*T {
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

	copy(result, q.data[q.current:m])
	if r > 0 {
		copy(result[m-q.current:], q.data[:r])
	}
	return result
}
