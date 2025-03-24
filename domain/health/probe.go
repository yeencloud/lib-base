package HealthProbe

type Status int

// Define constants using iota.
const (
	ProbeStatusUnhealthy Status = iota
	ProbeStatusDegraded
	ProbeStatusHealthy
)

func (s Status) String() string {
	return [...]string{"Unhealthy", "Degraded", "Healthy"}[s]
}

type HealthProbe interface {
	ProbeName() string
	ProbeStatus() Status
	ProbeDetails() []string
}
