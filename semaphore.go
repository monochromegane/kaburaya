package kaburaya

import (
	"math"
	"time"
)

// Semaphore represents an elastic semaphore.
type Semaphore struct {
	controller Controller
	reporter   Reporter
	ch         *elasticSemaphore
	done       chan struct{}
	Recorder   Recorder
}

// NewSem returns a semaphore.
func NewSem(duration time.Duration, gain float64) *Semaphore {
	sem := &Semaphore{
		controller: newDynamicTargetController(newPIDController(0.0, gain, gain, gain)),
		reporter:   newCPUReporter(),
		ch:         newElasticSemaphore(1),
		done:       make(chan struct{}),
	}
	go sem.adjust(duration)
	return sem
}

// Wait decrements semaphore. If semaphore will be negative,
// it blocks the process.
func (s *Semaphore) Wait() {
	s.ch.wait()
}

// Signal increments semaphore.
func (s *Semaphore) Signal() {
	s.ch.signal()
}

// Stop finalize resources.
func (s *Semaphore) Stop() {
	s.done <- struct{}{}
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
			l := s.ch.incrementLimit(round(inc))
			if s.Recorder != nil {
				s.Recorder.Record([]float64{usage, float64(l)})
			}
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
