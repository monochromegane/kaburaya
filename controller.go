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
