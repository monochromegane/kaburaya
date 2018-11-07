package kaburaya

import "sync"

type queue struct {
	counter int
	limit   int

	mu *sync.RWMutex
}

func newQueue(limit int) *queue {
	return &queue{
		limit: limit,
		mu:    &sync.RWMutex{},
	}
}

func (q *queue) incrementLimit(n int) int {
	defer q.mu.Unlock()
	q.mu.Lock()

	q.limit += n
	return q.limit
}

func (q *queue) full() bool {
	defer q.mu.RUnlock()
	q.mu.RLock()

	return q.counter == q.limit
}

func (q *queue) empty() bool {
	return q.counter == 0
}

func (q *queue) enqueue() {
	q.counter++
}

func (q *queue) dequeue() {
	q.counter--
}
