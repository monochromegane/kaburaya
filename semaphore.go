package kaburaya

import (
	"time"
)

// Semaphore represents an elastic semaphore.
type Semaphore struct {
	controller Controller
	ch         *elasticChannel
	done       chan struct{}
}

// NewSem returns a semaphore.
func NewSem(duration time.Duration) *Semaphore {
	sem := &Semaphore{
		controller: newQueueController(),
		ch:         newElasticChannel(1),
		done:       make(chan struct{}),
	}
	go sem.adjust(duration)
	return sem
}

// Wait decrements semaphore. If semaphore will be negative,
// it blocks the process.
func (s *Semaphore) Wait() {
	s.ch.send <- struct{}{}
}

// Signal increments semaphore.
func (s *Semaphore) Signal() {
	<-s.ch.receive
}

// Stop finalize resources.
func (s *Semaphore) Stop() {
	s.done <- struct{}{}
	s.ch.stop()
}

func (s *Semaphore) adjust(duration time.Duration) {
	t := time.NewTicker(duration)
	for {
		select {
		case <-t.C:
			inc := s.controller.Compute(float64(s.ch.queue.count()))
			s.ch.incrementLimit(int(inc))
		case <-s.done:
			t.Stop()
			break
		}
	}
}
