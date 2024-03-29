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
