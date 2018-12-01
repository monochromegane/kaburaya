package kaburaya

import (
	"sync"
)

type elasticSemaphore struct {
	limit   int
	counter int
	cond    *sync.Cond
}

func newElasticSemaphore(limit int) *elasticSemaphore {
	return &elasticSemaphore{
		limit:   limit,
		counter: limit,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (es *elasticSemaphore) wait() {
	es.cond.L.Lock()
	defer es.cond.L.Unlock()
	if es.counter == 0 {
		es.cond.Wait()
		return
	}
	es.counter--
	return
}

func (es *elasticSemaphore) signal() {
	es.cond.L.Lock()
	defer es.cond.L.Unlock()
	if es.counter == 0 {
		es.cond.Signal()
	}
	if es.limit >= es.counter+1 {
		es.counter++
	}
}

func (es *elasticSemaphore) incrementLimit(n int) int {
	es.cond.L.Lock()
	defer es.cond.L.Unlock()

	if n == 0 {
		return es.limit
	}

	if n > 0 {
		if es.counter == 0 {
			es.cond.Signal()
		}
		es.counter += n
		es.limit += n
	} else {
		newLimit := es.limit + n
		if newLimit < 1 {
			newLimit = 1
		}
		if newLimit < es.counter {
			es.counter = newLimit
		}
		es.limit = newLimit
	}
	return es.limit
}
