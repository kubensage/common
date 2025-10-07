package datastructure

import "sync"

// RingBuffer is a generic, thread-safe circular buffer for elements of type T.
//
// It has fixed capacity and uses FIFO semantics. When the buffer is full,
// inserting a new element overwrites the oldest one.
//
// All operations are safe for concurrent use by multiple goroutines.
type RingBuffer[T any] struct {
	data     []T        // underlying storage
	capacity int        // fixed capacity
	start    int        // index of the oldest element
	size     int        // current number of elements
	mu       sync.Mutex // mutex for thread safety
}

// NewRingBuffer creates a new empty RingBuffer with the given capacity.
//
// Parameters:
//   - cap: maximum number of elements the buffer can hold.
//
// Returns:
//
//	A pointer to a new RingBuffer[T].
func NewRingBuffer[T any](cap int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:     make([]T, cap),
		capacity: cap,
	}
}

// Add inserts an item into the buffer.
//
// If the buffer is not full, the item is added at the next free position.
// If the buffer is full, the oldest item is overwritten (circular behavior).
//
// Parameters:
//   - item: the value of type T to be added.
func (b *RingBuffer[T]) Add(item T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	idx := (b.start + b.size) % b.capacity
	b.data[idx] = item
	if b.size < b.capacity {
		b.size++
	} else {
		b.start = (b.start + 1) % b.capacity
	}
}

// Pop removes and returns the oldest item from the buffer.
//
// Returns:
//   - zero: the zero value of T (for compatibility; always returned).
//   - result: the oldest element in the buffer.
//   - ok: true if an element was successfully removed; false if buffer is empty.
func (b *RingBuffer[T]) Pop() (zero T, result T, ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.size == 0 {
		return zero, zero, false
	}
	result = b.data[b.start]
	var zeroValue T
	b.data[b.start] = zeroValue // remove reference for GC
	b.start = (b.start + 1) % b.capacity
	b.size--
	return zero, result, true
}

// Readd reinserts an item into the position it was last popped from,
// assuming space is available (i.e., the buffer is not full).
//
// This is useful when retrying failed processing of a popped item.
//
// Parameters:
//   - item: the item to reinsert.
//
// Returns:
//   - true if the item was successfully reinserted.
//   - false if the buffer is already full.
func (b *RingBuffer[T]) Readd(item T) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.size == b.capacity {
		return false
	}
	b.start = (b.start - 1 + b.capacity) % b.capacity
	b.data[b.start] = item
	b.size++
	return true
}

// Len returns the number of elements currently stored in the buffer.
//
// Returns:
//   - the current number of valid items (0 ≤ n ≤ capacity).
func (b *RingBuffer[T]) Len() int {
	return b.size
}
