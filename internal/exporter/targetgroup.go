package exporter

import (
	"strings"

	"github.com/postfinance/discovery"
)

// TargetGroup represents a prometheus target group.
type TargetGroup struct {
	Targets []string         `json:"targets,omitempty"`
	Labels  discovery.Labels `json:"labels,omitempty"`
}

// NewTargetGroup creates a new target group from discovery.Service.
func NewTargetGroup(s discovery.Service, cfg discovery.ExportConfig) TargetGroup {
	tg := TargetGroup{
		Labels: make(discovery.Labels),
	}

	if cfg == discovery.Disabled {
		return tg
	}

	if cfg == discovery.Blackbox {
		tg.Targets = []string{s.Endpoint.String()}
	} else {
		tg.Targets = []string{s.Endpoint.Host}
		tg.Labels["job"] = s.Name
		tg.Labels["instance"] = s.Endpoint.Host
		tg.Labels["__scheme__"] = s.Endpoint.Scheme
		tg.Labels["__metrics_path__"] = strings.TrimRight(s.Endpoint.Path, "/")
		for k, v := range s.Endpoint.Query() {
			tg.Labels["__param_"+k] = v[0]
		}
	}

	for k, v := range s.Labels {
		tg.Labels[k] = v
	}

	return tg
}
