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
	if q.limit < 1 {
		q.limit = 1
	}
	return q.limit
}

func (q *queue) full() bool {
	defer q.mu.RUnlock()
	q.mu.RLock()

	return q.counter >= q.limit
}

func (q *queue) empty() bool {
	defer q.mu.RUnlock()
	q.mu.RLock()

	return q.counter == 0
}

func (q *queue) enqueue() {
	defer q.mu.Unlock()
	q.mu.Lock()

	q.counter++
}

func (q *queue) dequeue() {
	defer q.mu.Unlock()
	q.mu.Lock()

	q.counter--
}

func (q *queue) count() int {
	defer q.mu.RUnlock()
	q.mu.RLock()

	return q.counter
}
