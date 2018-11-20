package kaburaya

import (
	"math"
	"time"
)

// Semaphore represents an elastic semaphore.
type Semaphore struct {
	controller Controller
	reporter   Reporter
	ch         *elasticChannel
	done       chan struct{}
}

// NewSem returns a semaphore.
func NewSem(duration time.Duration) *Semaphore {
	sem := &Semaphore{
		controller: newDynamicTargetController(newPIDController(0.0, 0.1, 0.5, 0.5)),
		reporter:   newCPUReporter(),
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
			usage, err := s.reporter.Report()
			if err != nil {
				break
			}
			inc := s.controller.Compute(usage)
			s.ch.incrementLimit(round(inc))
		case <-s.done:
			t.Stop()
			break
		}
	}
}

func round(n float64) int {
	if n > 0.0 {
		return int(math.Ceil(n))
	}
	return int(math.Floor(n))
}
