package kaburaya

const (
	modeInputEnable = iota + 1
	modeOutputEnable
)

type elasticChannel struct {
	send    chan<- struct{}
	receive <-chan struct{}
	queue   *queue
	done    chan<- struct{}
}

func newElasticChannel(limit int) *elasticChannel {
	in := make(chan struct{})
	out := make(chan struct{})
	queue := newQueue(limit)
	done := make(chan struct{})

	ec := &elasticChannel{
		send:    in,
		receive: out,
		queue:   queue,
		done:    done,
	}
	go poll(in, out, queue, done)
	return ec
}

func (ec *elasticChannel) incrementLimit(n int) int {
	return ec.queue.incrementLimit(n)
}

func (ec *elasticChannel) stop() {
	ec.done <- struct{}{}
	close(ec.send)
}

func poll(in <-chan struct{}, out chan<- struct{}, queue *queue, done <-chan struct{}) {
	state := modeInputEnable

loop:
	for {
		switch state {
		case modeInputEnable:
			select {
			case <-in:
				queue.enqueue()
				if queue.full() {
					state = modeOutputEnable
				} else {
					state = modeInputEnable | modeOutputEnable
				}
			case <-done:
				close(out)
				break loop
			}
		case modeInputEnable | modeOutputEnable:
			select {
			case <-in:
				queue.enqueue()
				if queue.full() {
					state = modeOutputEnable
				}
			case out <- struct{}{}:
				queue.dequeue()
				if queue.empty() {
					state = modeInputEnable
				}
			case <-done:
				close(out)
				break loop
			}
		case modeOutputEnable:
			select {
			case out <- struct{}{}:
				queue.dequeue()
				if queue.empty() {
					state = modeInputEnable
				} else if !queue.full() {
					state = modeInputEnable | modeOutputEnable
				}
			case <-done:
				close(out)
				break loop
			}
		}
	}
}
