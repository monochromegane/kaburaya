package kaburaya

import (
	"sync"
	"testing"
)

func TestElasticSemaphoreWaitAndSignalWithoutBlocking(t *testing.T) {
	limit := 2
	sem := newElasticSemaphore(limit)
	for i := 0; i < limit; i++ {
		sem.wait()
	}
	if sem.counter != 0 {
		t.Errorf("elasticSemaphore.wait should decrement the counter to 0, but %d.", sem.counter)
	}

	for i := 0; i < limit; i++ {
		sem.signal()
	}
	if sem.counter != 2 {
		t.Errorf("elasticSemaphore.signal should increment the counter to 2, but %d.", sem.counter)
	}
}

func TestElasticSemaphoreWaitAndSignalWithBlocking(t *testing.T) {
	limit := 2
	sem := newElasticSemaphore(limit)
	var wg sync.WaitGroup
	for i := 0; i < limit+1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer sem.signal()
			sem.wait()
		}()
	}
	wg.Wait()
	if sem.counter != 2 {
		t.Errorf("elasticSemaphore should increment the counter to 2, but %d.", sem.counter)
	}
}

func TestElasticSemaphoreIncrementLimitZero(t *testing.T) {
	limit := 2
	sem := newElasticSemaphore(limit)

	newLimit := sem.incrementLimit(0)
	if newLimit != limit {
		t.Errorf("elasticSemaphore.incrementLimit should not change the limit, but it changes to %d.", newLimit)
	}
}

func TestElasticSemaphoreIncrementLimitPositive(t *testing.T) {
	limit := 2
	sem := newElasticSemaphore(limit)

	inc := 1
	newLimit := sem.incrementLimit(inc)
	if newLimit != limit+inc {
		t.Errorf("elasticSemaphore.incrementLimit should change the limit to %d, but it changes to %d.", limit+inc, newLimit)
	}

	if sem.counter != limit+inc {
		t.Errorf("elasticSemaphore.incrementLimit should change the counter to %d, but it changes to %d.", limit+inc, sem.counter)
	}
}

func TestElasticSemaphoreIncrementLimitPositiveWithBlocking(t *testing.T) {
	limit := 2
	sem := newElasticSemaphore(limit)
	for i := 0; i < limit; i++ {
		sem.wait()
	}
	waiting := make(chan bool)
	done := make(chan bool)
	go func() {
		waiting <- true
		sem.wait()
		done <- true
	}()

	<-waiting // start blocking
	inc := 1
	newLimit := sem.incrementLimit(inc) // Increment limit and signal.
	<-done

	if newLimit != limit+inc {
		t.Errorf("elasticSemaphore.incrementLimit should change the limit to %d, but it changes to %d.", limit+inc, newLimit)
	}

	if sem.counter != 1 {
		t.Errorf("elasticSemaphore.incrementLimit should change the counter to %d, but it changes to %d.", 1, sem.counter)
	}
}

func TestElasticSemaphoreIncrementLimitNegative(t *testing.T) {
	limit := 2
	sem := newElasticSemaphore(limit)

	inc := -1
	newLimit := sem.incrementLimit(inc)
	if newLimit != limit+inc {
		t.Errorf("elasticSemaphore.incrementLimit should change the limit to %d, but it changes to %d.", limit+inc, newLimit)
	}

	if sem.counter != newLimit {
		t.Errorf("elasticSemaphore.incrementLimit should change the counter to %d, but it changes to %d.", newLimit, sem.counter)
	}
}

func TestElasticSemaphoreIncrementLimitNegativeUnderFlow(t *testing.T) {
	limit := 2
	sem := newElasticSemaphore(limit)

	inc := -2
	newLimit := sem.incrementLimit(inc)
	if newLimit != 1 {
		t.Errorf("elasticSemaphore.incrementLimit should change the limit to %d, but it changes to %d.", 1, newLimit)
	}

	if sem.counter != newLimit {
		t.Errorf("elasticSemaphore.incrementLimit should change the counter to %d, but it changes to %d.", newLimit, sem.counter)
	}
}
