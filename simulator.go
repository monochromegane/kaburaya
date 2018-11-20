package kaburaya

import (
	"fmt"
	"math"
)

type ControllerWithTarget interface {
	Compute(float64) float64
	CurrentTarget() float64
}

// Simulator represents a simulator of controller.
type Simulator struct {
	InitialNumWorker int
	Resource         float64
	Generator        generator
	Controller       ControllerWithTarget
	Reporter         Reporter
}

type Result struct {
	Usage     float64
	NumWorker int
	Target    float64
}

func (r Result) String() string {
	return fmt.Sprintf("%f,%d,%f\n", r.Usage, r.NumWorker, r.Target)
}

type FixController struct {
}

func (c *FixController) Compute(n float64) float64 {
	return 0.0
}

func (c *FixController) CurrentTarget() float64 {
	return 0.0
}

type SimpleController struct {
	Target float64
}

func (c *SimpleController) Compute(n float64) float64 {
	if n < c.Target {
		return 1.0
	}
	return 0.0
}

func (c *SimpleController) CurrentTarget() float64 {
	return c.Target
}

type PController struct {
	Target float64
	K      float64
}

func (c *PController) Compute(feedback float64) float64 {
	e := (c.Target - feedback)
	a := c.K * e
	// if a > 0.0 {
	// 	a = math.Ceil(a)
	// } else {
	// 	a = math.Floor(a)
	// }
	return a
}

func (c *PController) CurrentTarget() float64 {
	return c.Target
}

type PIController struct {
	Target float64
	Kp     float64
	Ki     float64
	I      float64
}

func (c *PIController) Compute(feedback float64) float64 {
	e := (c.Target - feedback)
	c.I += e
	// fmt.Printf("I: %f\n", c.I)
	return (c.Kp * e) + (c.Ki * c.I)
}

func (c *PIController) CurrentTarget() float64 {
	return c.Target
}

type PIShortController struct {
	Target  float64
	Kp      float64
	Ki      float64
	I       float64
	history []float64
	Span    int
}

func (c *PIShortController) Compute(feedback float64) float64 {
	e := (c.Target - feedback)

	if c.history == nil {
		c.history = make([]float64, c.Span)
	}
	c.history = append(c.history, e)[1:]

	c.I = c.sum()
	return (c.Kp * e) + (c.Ki * c.I)
}

func (c *PIShortController) Reset() {
	c.history = nil
}

func (c *PIShortController) sum() float64 {
	sum := 0.0
	for _, h := range c.history {
		sum += h
	}
	return sum
}

type PIDShortController struct {
	Target  float64
	Kp      float64
	Ki      float64
	Kd      float64
	I       float64
	D       float64
	history []float64
	Span    int
	reset   bool
}

func (c *PIDShortController) Compute(feedback float64) float64 {
	e := (c.Target - feedback)

	if c.history == nil {
		c.history = make([]float64, c.Span)
	}
	previous := c.history[len(c.history)-1]
	c.history = append(c.history, e)[1:]

	c.I = c.sum()
	if c.reset {
		c.I = 0.0
		c.reset = false
	}
	c.D = e - previous

	a := (c.Kp * e) + (c.Ki * c.I) + (c.Kd * c.D)
	if a > 0.0 {
		a = math.Ceil(a)
	} else {
		a = math.Floor(a)
	}
	return a
}

func (c *PIDShortController) Reset() {
	c.history = nil
	c.reset = true
}

func (c *PIDShortController) sum() float64 {
	sum := 0.0
	for _, h := range c.history {
		sum += h
	}
	return sum
}

type DynamicController struct {
	// Controller *PController
	// Controller *PIController
	Controller *PIDShortController
	Span       int
	counter    int
	history    []float64
	ema        *EMA
	previous   float64
}

func (c *DynamicController) Compute(feedback float64) float64 {
	previous := 0.01
	if c.previous != 0.0 {
		previous = c.previous
	}
	changeRate := (feedback - previous) / previous
	c.previous = feedback

	if c.ema == nil {
		c.ema = &EMA{Span: c.Span}
	}
	if c.history == nil {
		c.history = make([]float64, c.Span)
	}
	// First edition
	// c.history = append(c.history, feedback)[1:]
	// if c.counter != 0 && c.counter%c.Span == 0 {
	// 	c.Controller.Target = avg(c.history) - 10.0
	// }

	// OK edition
	// c.history = append(c.history, feedback)[1:]
	// if changeRate > 0.3 {
	// 	c.Controller.Reset()
	// 	c.Controller.Target = feedback + 2.0
	// 	// c.Controller.Target = (feedback + c.Controller.Target) / 2.0
	// 	// c.Controller.Target = 100.0
	// 	// target := feedback * 3.0
	// 	// if target >= 100.0 {
	// 	// 	target = 100.0
	// 	// }
	// 	// c.Controller.Target = target
	// } else if changeRate == 0.0 {
	// 	// c.Controller.Reset()
	// 	if feedback == 0 {
	// 		c.Controller.Target = -0.0
	// 	} else {
	// 		c.Controller.Target = feedback - 2.0
	// 	}
	// } else if changeRate < -0.3 {
	// 	c.Controller.Reset()
	// 	if feedback == 0 {
	// 		c.Controller.Target = -100.0
	// 	} else {
	// 		// c.Controller.Target = feedback // + 1.0
	// 		// c.Controller.Target = feedback / 2.0
	// 		// target := c.Controller.Target
	// 		// c.Controller.Target = (feedback + c.Controller.Target) / 2.0
	// 		c.Controller.Target = feedback + 2.0
	// 		// fmt.Printf("(%f + %f)/2.0 = %f\n", feedback, target, c.Controller.Target)
	// 	}
	// }

	// OK edition (step by step)
	c.history = append(c.history, feedback)[1:]
	if changeRate > 0.3 {
		c.Controller.Reset()
		c.Controller.Target = feedback + 2.0
		// c.Controller.Target = (feedback + c.Controller.Target) / 2.0
		// c.Controller.Target = 100.0
		// target := feedback * 3.0
		// if target >= 100.0 {
		// 	target = 100.0
		// }
		// c.Controller.Target = target
	} else if changeRate == 0.0 {
		c.Controller.Reset()
		if feedback == 0 {
			c.Controller.Target = -0.0
		} else {
			c.Controller.Target = feedback - 2.0
		}
	} else if changeRate < -0.3 {
		c.Controller.Reset()
		if feedback == 0 {
			c.Controller.Target = -100.0
		} else {
			// c.Controller.Target = feedback // + 1.0
			// c.Controller.Target = feedback / 2.0
			// target := c.Controller.Target
			// c.Controller.Target = (feedback + c.Controller.Target) / 2.0
			c.Controller.Target = feedback + 2.0
			// c.Controller.Target = -100.0
			// fmt.Printf("(%f + %f)/2.0 = %f\n", feedback, target, c.Controller.Target)
		}
	}

	// Use Avg
	// c.history = append(c.history, feedback)[1:]
	// if c.counter > c.Span && changeRate > 0.3 {
	// 	c.Controller.Target = 100.0
	// } else if c.counter > c.Span && (changeRate == 0.0 || changeRate < -0.3) {
	// 	if feedback == 0 {
	// 		c.Controller.Target = -100.0
	// 	} else {
	// 		c.Controller.Target = feedback // + 1.0
	// 	}
	// }

	// } else if c.counter != 0 && c.counter%c.Span == 0 {
	// 	if avg := avg(c.history); avg == 0.0 {
	// 		c.Controller.Target = avg - 10.0
	// 	} else {
	// 		c.Controller.Target = avg
	// 	}
	// }
	// } else if (c.counter > c.Span && changeRate == 0.0) || (c.counter > c.Span && math.Abs(changeRate) > 0.3) {
	// 	if feedback == 0 {
	// 		c.Controller.Target = -100.0
	// 	} else {
	// 		c.Controller.Target = feedback // + 1.0
	// 	}
	// 	// c.Controller.Target = feedback // + 1.0
	// 	// if c.Controller.Target == 0.0 {
	// 	// 	c.Controller.Target = -100.0
	// 	// }
	// }

	// Use EMA
	// avg := c.ema.Avg(feedback)
	// if (c.counter != 0 && c.counter%c.Span == 0) || (c.counter > c.Span && math.Abs(changeRate) > 0.5) {
	// 	if avg == 0.0 {
	// 		c.Controller.Target = avg - 100.0
	// 	} else {
	// 		c.Controller.Target = avg
	// 	}
	// } else if c.counter > c.Span && changeRate == 0.0 {
	// 	c.Controller.Target = feedback
	// }

	c.counter++
	a := c.Controller.Compute(feedback)
	// if c.counter != 1 && feedback == 0.0 && a > 0.0 {
	// 	a = -10.0
	// }
	fmt.Printf("%03d Target: %f, Feedback: %f, Compute: %f, ChangeRate: %f, I: %f\n", c.counter, c.Controller.Target, feedback, a, changeRate, c.Controller.I)
	return a
}

func (c *DynamicController) CurrentTarget() float64 {
	return c.Controller.Target
}

type EMA struct {
	Span     int
	previous float64
	history  []float64
}

func (e *EMA) Avg(x float64) float64 {
	if e.history == nil {
		e.history = []float64{}
	}
	e.history = append(e.history, x)
	if len(e.history) < e.Span {
		return avg(e.history)
	}

	alpha := 2.0 / (float64(e.Span) + 1.0)
	if len(e.history) == e.Span {
		avg := avg(e.history)
		ema := avg + alpha*(x-avg)
		e.previous = ema
	} else {
		ema := e.previous + alpha*(x-e.previous)
		e.previous = ema
	}
	return e.previous
}

type StabilityController struct {
	Controller *PController
	Span       int
	SD         float64
	history    []float64
}

func (c *StabilityController) Compute(feedback float64) float64 {
	if c.history == nil {
		c.history = make([]float64, c.Span)
	}
	history := append(c.history, feedback)[1:]
	avg := avg(history)
	sum := 0.0
	for _, h := range history {
		sum += math.Pow(h-avg, 2)
	}
	v := sum / float64(len(history))
	sd := math.Sqrt(v)
	c.history = history
	fmt.Printf("SD: %f\n", sd)
	if sd > c.SD {
		return c.Controller.Compute(feedback)
	}
	return -1.0

}

func (c *StabilityController) CurrentTarget() float64 {
	return c.Controller.Target
}

type RateController struct {
	previous   float64
	Controller *PController
}

func (c *RateController) Compute(feedback float64) float64 {
	previous := 0.01
	if c.previous != 0.0 {
		previous = c.previous
	}
	changeRate := (feedback - previous) / previous
	c.previous = feedback

	a := -1.0 * c.Controller.Compute(changeRate)
	fmt.Printf("Feedback: %f, ChangeRate: %f, Controller: %f\n", feedback, changeRate, a)
	return a
}

func (c *RateController) CurrentTarget() float64 {
	return c.Controller.Target
}

type DynamicRateController struct {
	Controller *RateController
	Span       int
	counter    int
	history    []float64
	previous   float64
}

func (c *DynamicRateController) Compute(feedback float64) float64 {
	previous := 0.01
	if c.previous != 0.0 {
		previous = c.previous
	}
	changeRate := (feedback - previous) / previous
	c.previous = feedback

	if c.history == nil {
		c.history = make([]float64, c.Span)
	}
	c.history = append(c.history, changeRate)[1:]
	if c.counter != 0 && c.counter%c.Span == 0 {
		c.Controller.Controller.Target = avg(c.history)
	}
	fmt.Printf("%f\n", c.Controller.Controller.Target)
	c.counter++
	return c.Controller.Compute(feedback)
}

func (c *DynamicRateController) CurrentTarget() float64 {
	return c.Controller.Controller.Target
}

// Run runs simulation.
func (s *Simulator) Run(step int) []Result {
	results := make([]Result, step+1)
	results[0] = Result{Usage: 0.0, NumWorker: s.InitialNumWorker, Target: s.Controller.CurrentTarget()}
	numWorker := s.InitialNumWorker
	jobs := newJobs()

	total := 0
	beforeUsage := 0.0
	for i := 0; i < step; i++ {
		// Compute worker
		numWorker += int(s.Controller.Compute(beforeUsage))
		if numWorker <= 0 {
			numWorker = 1
		}
		total += numWorker

		// Push new jobs
		newWorkloads := s.Generator.Generate(i)
		for _, workload := range newWorkloads {
			jobs.add(workload)
		}

		// Do jobs
		jobs.reset()
		remain := s.Resource
		for j := 0; j < numWorker; j++ {
			job := jobs.processable()
			if job == nil {
				break
			}
			usage := job.process(remain)
			remain -= usage
			if remain <= 0.0 {
				break
			}
		}

		// Process for 0 workload
		jobs.wait()

		// if jobs.isFinished() {
		// 	fmt.Printf("Break: %d -> Total: %d\n", i+1, total)
		// }

		// result
		results[i+1] = Result{Usage: s.Resource - remain, NumWorker: numWorker, Target: s.Controller.CurrentTarget()}

		beforeUsage = results[i+1].Usage
	}

	return results
}

type jobs struct {
	jobs []*job
}

func newJobs() *jobs {
	return &jobs{jobs: []*job{}}
}

func (js *jobs) add(workloads []float64) {
	js.jobs = append(js.jobs, &job{workloads: workloads})
}

func (js *jobs) reset() {
	for _, job := range js.jobs {
		job.processed = false
	}
}

func (js jobs) processable() *job {
	for _, job := range js.jobs {
		if job.isFinished() || job.processed || job.isWaiting() {
			continue
		}
		return job
	}
	return nil
}

func (js jobs) wait() {
	for _, job := range js.jobs {
		if !job.isFinished() && !job.processed && !job.isWaiting() {
			job.process(0.0)
		}
	}
}

func (js jobs) isFinished() bool {
	for _, job := range js.jobs {
		if !job.isFinished() {
			return false
		}
	}
	return true
}

func (js jobs) numQueue() int {
	n := 0
	for _, job := range js.jobs {
		if job.isFinished() {
			continue
		}
		n++
	}
	return n
}

type job struct {
	workloads []float64
	processed bool
}

func (j *job) process(work float64) float64 {
	workloads := j.workloads[0]
	j.processed = true
	if workloads <= work {
		j.workloads = j.workloads[1:]
		return workloads
	}
	j.workloads[0] = workloads - work
	return work
}

func (j job) isWaiting() bool {
	return j.workloads[0] == 0
}

func (j job) isFinished() bool {
	return len(j.workloads) == 0
}

type generator interface {
	Generate(int) [][]float64
}
