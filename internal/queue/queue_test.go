package queue

import (
	"testing"
)

func TestQueue_RemainingFilled(t *testing.T) {
	runes := "abcdefghijklmnopqrstuvwxyz"
	capacity := 20
	q := New[rune](capacity)
	for _, r := range runes {
		err := q.Push(r)
		if err == ErrCapacityExceeded {
			_, _ = q.Pop()
			_ = q.Push(r)
		}
	}
	remaining := q.Remaining()
	if len(remaining) != q.size {
		t.Errorf("Expected %d remaining, got %d", q.size, len(remaining))
	}
	for i, r := range runes[len(runes)-capacity:] {
		if r != *(remaining[i]) {
			t.Errorf("Expected(at %d) %#v, got %#v", i, string(*remaining[i]), string(r))
		}
	}
}

func TestQueue_RemainingSinglePass(t *testing.T) {
	runes := "abcdefghijklmnopqrstuvwxyz"
	capacity, rem := 20, 15
	q := New[rune](capacity)
	for _, r := range runes {
		err := q.Push(r)
		if err == ErrCapacityExceeded {
			_, _ = q.Pop()
			_ = q.Push(r)
		}
	}
	for i := 0; i < rem; i++ {
		_, _ = q.Pop()
	}
	remaining := q.Remaining()
	if len(remaining) != q.size {
		t.Errorf("Expected %d remaining, got %d", q.size, len(remaining))
	}
	for i, r := range runes[len(runes)-capacity+rem:] {
		if r != *(remaining[i]) {
			t.Errorf("Expected(at %d) %#v, got %#v", i, string(*remaining[i]), string(r))
		}
	}
}
