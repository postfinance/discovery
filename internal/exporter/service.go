package exporter

import (
	"sort"

	"github.com/postfinance/discovery"
)

type service struct {
	discovery.Service
}

func (s service) key() string {
	return s.Namespace + ":" + s.Name
}

// isBlackbox returns true if service has "blackbox" as exportConfigMetaLabel label.
// func (s service) isBlackbox() bool {
// return s.Labels.Get(exportConfigMetaLabel) == blackboxName
// }

type services map[string]service

func (s services) list() []service {
	l := make([]service, 0, len(s))

	for i := range s {
		l = append(l, s[i])
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].Endpoint.String() < l[j].Endpoint.String()
	})

	return l
}
