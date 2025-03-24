package health

import (
	"github.com/yeencloud/lib-base/domain/health"
)

type Probe struct {
	Hostname string
	Probes   []HealthProbe.HealthProbe `json:"probes"`
}

type ServiceHealth struct {
	Hostname string `json:"hostname"`

	RawStatus HealthProbe.Status `json:"rawStatus"`
	Status    string             `json:"status"`
}

func (p *Probe) Health() *ServiceHealth {
	currentStatus := HealthProbe.ProbeStatusHealthy
	for _, probe := range p.Probes {
		probeStatus := probe.ProbeStatus()
		currentStatus = min(currentStatus, probeStatus)
	}

	return &ServiceHealth{
		RawStatus: currentStatus,
		Status:    currentStatus.String(),
		Hostname:  p.Hostname,
	}
}

func NewHealthProbe(hostname string) *Probe {
	return &Probe{
		Hostname: hostname,
	}
}
