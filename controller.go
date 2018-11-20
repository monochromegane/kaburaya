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
