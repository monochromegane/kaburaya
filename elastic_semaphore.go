package kaburaya

import (
	"sync"
)

type elasticSemaphore struct {
	limit      int
	counter    int
	mu         *sync.Mutex
	waitChan   chan struct{}
	signalChan chan struct{}
	waiting    bool
}

func newElasticSemaphore(limit int) *elasticSemaphore {
	return &elasticSemaphore{
		limit:      limit,
		counter:    limit,
		mu:         &sync.Mutex{},
		waitChan:   make(chan struct{}),
		signalChan: make(chan struct{}),
		waiting:    false,
	}
}

func (es *elasticSemaphore) wait() {
WAIT:
	es.mu.Lock()
	if es.counter <= 0 {
		es.counter--
		es.waiting = true
		es.mu.Unlock()
		es.waitChan <- struct{}{}
		<-es.signalChan
		es.mu.Lock()
		if es.counter < 0 {
			es.mu.Unlock()
			goto WAIT
		}
		es.mu.Unlock()
		return
	}
	es.waiting = false
	es.counter--
	es.mu.Unlock()
	return
}

func (es *elasticSemaphore) signal() {
WAIT:
	es.mu.Lock()
	doSignal := false
	if es.waiting && es.counter == -1 {
		select {
		case <-es.waitChan:
		default:
			es.mu.Unlock()
			goto WAIT
		}
		doSignal = true
	}
	if es.limit >= es.counter+1 {
		es.counter++
	} else {
		es.counter = es.limit
	}
	if doSignal {
		es.signalChan <- struct{}{}
	}
	es.mu.Unlock()
}

func (es *elasticSemaphore) incrementLimit(n int) int {
	if n == 0 {
		return es.limit
	}

WAIT:
	es.mu.Lock()
	if n > 0 {
		doSignal := false
		if es.waiting && (es.counter == -1 || (es.counter < 0 && es.counter+n > 0)) {
			select {
			case <-es.waitChan:
			default:
				es.mu.Unlock()
				goto WAIT
			}
			doSignal = true
		}
		es.counter += n
		es.limit += n
		if doSignal {
			es.signalChan <- struct{}{}
		}
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
	defer es.mu.Unlock()
	return es.limit
}
