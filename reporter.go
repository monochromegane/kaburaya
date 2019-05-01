package kaburaya

import (
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
	if u == 0 && n == 0 && s == 0 && i == 0 {
		return 0.0
	}

	return (float64(u+n+s) / float64(u+n+s+i)) * 100
}
