package kaburaya

// Semaphore represents an elastic semaphore.
type Semaphore struct {
	ch *elasticChannel
}

// NewSem returns a semaphore.
func NewSem() *Semaphore {
	return &Semaphore{
		ch: newElasticChannel(1),
	}
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
	s.ch.stop()
}
