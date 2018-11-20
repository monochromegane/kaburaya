package kaburaya

import (
	"math"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

// Reporter represents an interface of resource reporter.
type Reporter interface {
	Report() (float64, error)
}

const (
	statFile = "/proc/stat"
)

type cpuReporter struct {
	previous *linuxproc.Stat
}

func newCPUReporter() *cpuReporter {
	return &cpuReporter{}
}

func (r *cpuReporter) Report() (float64, error) {
	return r.usage()
}

func (r *cpuReporter) usage() (float64, error) {
	current, err := stats()
	if err != nil {
		return 0.0, err
	}
	usage := usage(r.previous, current)
	r.previous = current
	return usage, nil
}

type cpuStabilityReporter struct {
	previous *linuxproc.Stat
	span     int
	history  []float64
}

func newCPUStabilityReporter(span int) *cpuStabilityReporter {
	return &cpuStabilityReporter{
		span:    span,
		history: make([]float64, span),
	}
}

func (r *cpuStabilityReporter) Report() (float64, error) {
	usage, err := r.usage()
	if err != nil {
		return usage, err
	}
	history := append(r.history, usage)[1:]
	avg := avg(history)
	sum := 0.0
	for _, h := range history {
		sum += math.Pow(h-avg, 2)
	}
	v := sum / float64(len(history))
	sd := math.Sqrt(v)
	r.history = history

	return sd, nil
}

func (r *cpuStabilityReporter) usage() (float64, error) {
	current, err := stats()
	if err != nil {
		return 0.0, err
	}
	usage := usage(r.previous, current)
	r.previous = current
	return usage, nil
}

func stats() (*linuxproc.Stat, error) {
	return linuxproc.ReadStat(statFile)
}

func usage(previous, current *linuxproc.Stat) float64 {
	if previous == nil {
		return 0.0
	}
	u := current.CPUStatAll.User - previous.CPUStatAll.User
	n := current.CPUStatAll.Nice - previous.CPUStatAll.Nice
	s := current.CPUStatAll.System - previous.CPUStatAll.System
	i := current.CPUStatAll.Idle - previous.CPUStatAll.Idle

	return (float64(u+n+s) / float64(u+n+s+i)) * 100
}

func avg(ms []float64) float64 {
	total := 0.0
	for _, m := range ms {
		total += m
	}
	return total / float64(len(ms))
}
