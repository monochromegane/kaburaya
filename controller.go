package kaburaya

// Controller represents an interface of controller for elastic semaphore.
type Controller interface {
	Compute(float64) float64
}

func newQueueController() *queueController {
	return &queueController{previous: 0}
}

type queueController struct {
	previous int
}

func (c *queueController) Compute(current float64) float64 {
	// TODO: Be smart
	cur := int(current)
	if cur >= c.previous {
		c.previous = cur
		return 1.0
	}
	c.previous = cur
	return -1.0
}

type pidController struct {
	target   float64
	kp       float64
	ki       float64
	kd       float64
	i        float64
	previous float64
}

func newPIDController(target, kp, ki, kd float64) *pidController {
	return &pidController{
		target: target,
		kp:     kp,
		ki:     ki,
		kd:     kd,
	}
}

func (c *pidController) Compute(feedback float64) float64 {
	e := (c.target - feedback)
	c.i += e
	d := e - c.previous
	c.previous = e

	return (c.kp * e) + (c.ki * c.i) + (c.kd * d)
}

type dynamicTargetController struct {
	controller *pidController
	previous   float64
	changeRate float64
}

func newDynamicTargetController(controller *pidController) *dynamicTargetController {
	return &dynamicTargetController{
		controller: controller,
		previous:   0.01,
		changeRate: 0.3,
	}
}

const defaultPrevious = 0.01

func (c *dynamicTargetController) Compute(feedback float64) float64 {
	previous := c.previous
	if previous == 0.0 {
		previous = defaultPrevious // Avoid zero div
	}
	changeRate := (feedback - previous) / previous
	c.previous = feedback

	if changeRate == 0.0 {
		c.controller.i = 0.0
		if feedback != 0.0 {
			c.controller.target = feedback - 2.0
		}
	} else if changeRate > c.changeRate || changeRate < -c.changeRate {
		c.controller.i = 0.0
		if feedback == 0.0 {
			c.controller.target = -100.0
		} else {
			c.controller.target = feedback + 2.0
		}
	}

	return c.controller.Compute(feedback)
}
