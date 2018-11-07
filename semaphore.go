package kaburaya

// Semaphore represents an elastic semaphore.
type Semaphore struct {
}

// NewSem returns a semaphore.
func NewSem() *Semaphore {
	return &Semaphore{}
}

// Wait decrements semaphore. If semaphore will be negative,
// it blocks the process.
func (s *Semaphore) Wait() {
}

// Signal increments semaphore.
func (s *Semaphore) Signal() {
}
