package kaburaya

import (
	"testing"
)

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
