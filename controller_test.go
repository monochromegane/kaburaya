package kaburaya

import "testing"

func TestQueueControllerCompute(t *testing.T) {
	c := newQueueController()
	inc := c.Compute(1.0)
	if inc != 1.0 {
		t.Errorf("queueController.Compute should return incremental value when current > previous, but %f", inc)
	}

	inc = c.Compute(1.0)
	if inc != 0.0 {
		t.Errorf("queueController.Compute should return 0.0 when current == previous, but %f", inc)
	}

	inc = c.Compute(0.0)
	if inc != 0.0 {
		t.Errorf("queueController.Compute should return 0.0 when current < previous, but %f", inc)
	}
}
