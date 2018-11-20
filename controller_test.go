package kaburaya

import (
	"testing"
)

func TestQueueControllerCompute(t *testing.T) {
	c := newQueueController()
	inc := c.Compute(1.0)
	if inc != 1.0 {
		t.Errorf("queueController.Compute should return incremental value when current > previous, but %f", inc)
	}

	inc = c.Compute(1.0)
	if inc != 1.0 {
		t.Errorf("queueController.Compute should return 1.0 when current == previous, but %f", inc)
	}

	inc = c.Compute(0.0)
	if inc != -1.0 {
		t.Errorf("queueController.Compute should return -1.0 when current < previous, but %f", inc)
	}
}

func TestPIDControllerCompute(t *testing.T) {
	c := newPIDController(100.0, 0.1, 0.1, 0.1)
	expects := []float64{30.0, 21.0, 33.7}
	feedback := 0.0
	for i := 0; i < 3; i++ {
		feedback = c.Compute(feedback)
		if feedback != expects[i] {
			t.Errorf("pidController.Compute should return %f, but %f", expects[i], feedback)
		}
	}
}

func TestDynamicTargetControllerCompute(t *testing.T) {
	var pid *pidController
	var c *dynamicTargetController
	var feedback float64

	// ChangeRate == 0.0
	pid = newPIDController(0.0, 0.1, 0.1, 0.1)
	c = newDynamicTargetController(pid)
	feedback = defaultPrevious
	c.Compute(feedback)
	if expect := feedback - 2.0; c.controller.target != expect {
		t.Errorf("dynamicTargetController.Compute should change target to %f, but %f", expect, c.controller.target)
	}

	// ChangeRate > 0.3
	pid = newPIDController(0.0, 0.1, 0.1, 0.1)
	c = newDynamicTargetController(pid)
	feedback = 1.0
	c.Compute(feedback)
	if expect := feedback + 2.0; c.controller.target != expect {
		t.Errorf("dynamicTargetController.Compute should change target to %f, but %f", expect, c.controller.target)
	}

	// ChangeRate < -0.3 && feedback > 0.0
	pid = newPIDController(0.0, 0.1, 0.1, 0.1)
	c = newDynamicTargetController(pid)
	feedback = 1.0
	c.Compute(feedback)
	if expect := feedback + 2.0; c.controller.target != expect {
		t.Errorf("dynamicTargetController.Compute should change target to %f, but %f", expect, c.controller.target)
	}

	// ChangeRate < -0.3 && feedback == 0.0
	pid = newPIDController(0.0, 0.1, 0.1, 0.1)
	c = newDynamicTargetController(pid)
	feedback = 0.0
	c.Compute(feedback)
	if expect := -100.0; c.controller.target != expect {
		t.Errorf("dynamicTargetController.Compute should change target to %f, but %f", expect, c.controller.target)
	}
}
