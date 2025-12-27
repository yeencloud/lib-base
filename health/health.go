package health

import (
	"github.com/yeencloud/lib-base/domain/health"
	"github.com/yeencloud/lib-shared/env"
)

type Probe struct {
	hostname string
	build    env.Build
}

type ServiceHealth struct {
	Hostname string `json:"hostname"`

	Status    string             `json:"status"`
	RawStatus HealthProbe.Status `json:"rawStatus"`

	Repo    string `json:"repo"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func (p *Probe) Health() *ServiceHealth {
	currentStatus := HealthProbe.ProbeStatusHealthy

	return &ServiceHealth{
		Hostname: p.hostname,

		Status:    currentStatus.String(),
		RawStatus: currentStatus,

		Repo:    p.build.Repository,
		Version: p.build.AppVersion,
		Commit:  p.build.Commit,
	}
}

func NewHealthProbe(hostname string, build env.Build) *Probe {
	return &Probe{
		hostname: hostname,
		build:    build,
	}
}
