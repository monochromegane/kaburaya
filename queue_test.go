package kaburaya

import "testing"

func TestQueueEnqueueAndDequeue(t *testing.T) {
	queue := newQueue(1)
	queue.enqueue()
	if queue.counter != 1 {
		t.Errorf("queue.enqueue should increment counter to 1, but %d", queue.counter)
	}

	queue.dequeue()
	if queue.counter != 0 {
		t.Errorf("queue.dequeue should decrement counter to 0, but %d", queue.counter)
	}
}

func TestQueueFullAndEmpty(t *testing.T) {
	queue := newQueue(1)
	queue.enqueue()
	if !queue.full() {
		t.Errorf("queue.full should return true when queue is full, but false")
	}
	queue.dequeue()
	if !queue.empty() {
		t.Errorf("queue.empty should return true when queue is empty, but false")
	}
}

func TestQueueIncrementLimit(t *testing.T) {
	queue := newQueue(1)
	queue.enqueue()
	if !queue.full() {
		t.Errorf("queue.full should return true when queue is full, but false")
	}

	newLimit := queue.incrementLimit(1)
	if newLimit != 2 {
		t.Errorf("queue.incrementLimit should return incremented limit value, but %d", newLimit)
	}
	if queue.full() {
		t.Errorf("queue.incrementLimit should increment limit")
	}
}
