package kaburaya

import (
	"testing"
)

func TestElasticChannel(t *testing.T) {
	ec := newElasticChannel(2)
	defer ec.stop()

	for i := 0; i < 2; i++ {
		ec.send <- struct{}{}
	}

	select {
	case ec.send <- struct{}{}:
		t.Errorf("Elastic channel with 2 buffer can 2 time send.")
	default:
	}

	for i := 0; i < 2; i++ {
		<-ec.receive
	}

	select {
	case <-ec.receive:
		t.Errorf("Elastic channel with 2 buffer can only 2 time receive.")
	default:
	}
}
