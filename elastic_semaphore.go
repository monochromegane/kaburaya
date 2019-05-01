package kaburaya

import (
	"sync"
)

type elasticSemaphore struct {
	limit   int
	counter int
	cond    *sync.Cond
	mu      *sync.Mutex
}

func newElasticSemaphore(limit int) *elasticSemaphore {
	return &elasticSemaphore{
		limit:   limit,
		counter: limit,
		cond:    sync.NewCond(&sync.Mutex{}),
		mu:      &sync.Mutex{},
	}
}

func (es *elasticSemaphore) wait() {
	es.cond.L.Lock()
	defer es.cond.L.Unlock()
WAIT:
	es.mu.Lock()
	if es.counter <= 0 {
		es.mu.Unlock()
		es.cond.Wait()
		es.mu.Lock()
		if es.counter > 0 {
			es.counter--
		} else {
			es.mu.Unlock()
			goto WAIT
		}
		es.mu.Unlock()
		return
	}
	es.counter--
	es.mu.Unlock()
	return
}

func (es *elasticSemaphore) signal() {
	es.mu.Lock()
	defer es.mu.Unlock()
	if es.counter == 0 {
		es.cond.Signal()
	}
	if es.limit >= es.counter+1 {
		es.counter++
	} else {
		es.counter = es.limit
	}
}

func (es *elasticSemaphore) incrementLimit(n int) int {
	es.mu.Lock()
	defer es.mu.Unlock()

	if n == 0 {
		return es.limit
	}

	if n > 0 {
		if es.counter == 0 || (es.counter < 0 && es.counter+n > 0) {
			es.cond.Signal()
		}
		es.counter += n
		es.limit += n
	} else {
		newLimit := es.limit + n
		if newLimit < 1 {
			newLimit = 1
		}
		if es.limit != newLimit {
			es.counter = newLimit - (es.limit - es.counter)
		}
		es.limit = newLimit
	}
	return es.limit
}
