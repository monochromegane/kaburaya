package kaburaya

// Recorder represents recorder for semaphore metrics.
type Recorder interface {
	Record([]float64)
}
