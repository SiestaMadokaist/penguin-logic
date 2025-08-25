package queue

type qu[T any] struct {
	entries []T
	head    int
}

type Queue[T any] interface {
	Enqueue(item T)
	Peek() T
	Dequeue() (bool, T)
}

func New[T any](initial []T) Queue[T] {
	return &qu[T]{entries: initial, head: 0}
}

func (q *qu[T]) Enqueue(item T) {
	q.entries = append(q.entries, item)
}

func (q *qu[T]) Peek() T {
	return q.entries[q.head]
}

func (q *qu[T]) Dequeue() (bool, T) {
	if (q.head) >= len(q.entries) {
		var zero T
		return false, zero
	}
	item := q.entries[q.head]
	q.head++
	return true, item
}
