// Package convert converts proto types to domain types and
// vice versa.
package convert

import (
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/zbindenren/discovery"
	discoveryv1 "github.com/zbindenren/discovery/pkg/discoverypb"
)

// NamespaceToPB converts *discovery.Namespace to *discoveryv1.Namespace.
func NamespaceToPB(n *discovery.Namespace) *discoveryv1.Namespace {
	pb := &discoveryv1.Namespace{
		Name:     n.Name,
		Export:   int32(n.Export),
		Modified: timeToPB(&n.Modified),
	}

	return pb
}

// NamespaceFromPB converts *discovery.Namespace to *discoveryv1.Namespace.
func NamespaceFromPB(pb *discoveryv1.Namespace) *discovery.Namespace {
	n := &discovery.Namespace{
		Name:     pb.Name,
		Export:   discovery.ExportConfig(pb.Export),
		Modified: timeFromPB(pb.Modified),
	}

	return n
}

// ServerToPB converts *discovery.Server to *discoveryv1.Server.
func ServerToPB(s *discovery.Server) *discoveryv1.Server {
	pb := &discoveryv1.Server{
		Name:     s.Name,
		Labels:   s.Labels,
		Modified: timeToPB(&s.Modified),
	}

	return pb
}

// ServerFromPB converts *discovery.Server to *discoveryv1.Server.
func ServerFromPB(pb *discoveryv1.Server) *discovery.Server {
	s := &discovery.Server{
		Name:     pb.Name,
		Labels:   pb.Labels,
		Modified: timeFromPB(pb.Modified),
	}

	return s
}

// ServersToPB converts discovery.Servers to slices of *discoveryv1.Server.
func ServersToPB(s discovery.Servers) []*discoveryv1.Server {
	result := make([]*discoveryv1.Server, 0, len(s))

	for i := range s {
		result = append(result, ServerToPB(&s[i]))
	}

	return result
}

// ServersFromPB converts slices *discoveryv1.Server discovery.Servers.
func ServersFromPB(s []*discoveryv1.Server) discovery.Servers {
	result := make(discovery.Servers, 0, len(s))

	for i := range s {
		result = append(result, *ServerFromPB(s[i]))
	}

	return result
}

// ServiceToPB converts *discovery.Service to *discoveryv1.Service.
func ServiceToPB(s *discovery.Service) *discoveryv1.Service {
	pb := &discoveryv1.Service{
		Id:          s.ID,
		Name:        s.Name,
		Labels:      s.Labels,
		Description: s.Description,
		Endpoint:    s.Endpoint.String(),
		Namespace:   s.Namespace,
		Selector:    s.Selector,
		Servers:     s.Servers,
		Modified:    timeToPB(&s.Modified),
	}

	return pb
}

// ServiceFromPB converts *discoveryv1.Service *discovery.Service.
func ServiceFromPB(pb *discoveryv1.Service) *discovery.Service {
	e, _ := url.Parse(pb.GetEndpoint())
	s := &discovery.Service{
		ID:          pb.GetId(),
		Name:        pb.GetName(),
		Labels:      pb.GetLabels(),
		Description: pb.GetDescription(),
		Endpoint:    e,
		Namespace:   pb.GetNamespace(),
		Selector:    pb.GetSelector(),
		Servers:     pb.GetServers(),
		Modified:    timeFromPB(pb.GetModified()),
	}

	return s
}

// ServicesToPB converts discovery.Services to slice of *discoveryv1.Service.
func ServicesToPB(s discovery.Services) []*discoveryv1.Service {
	result := make([]*discoveryv1.Service, 0, len(s))

	for i := range s {
		result = append(result, ServiceToPB(&s[i]))
	}

	return result
}

// ServicesFromPB converts slice of *discoveryv1.Service to discovery.Services.
func ServicesFromPB(s []*discoveryv1.Service) discovery.Services {
	result := make(discovery.Services, 0, len(s))

	for i := range s {
		result = append(result, *ServiceFromPB(s[i]))
	}

	return result
}

func timeFromPB(ts *timestamp.Timestamp) time.Time {
	var t time.Time

	if ts == nil {
		t = time.Unix(0, 0).UTC() // treat nil like the empty Timestamp
	} else {
		t = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	}

	return t
}

func timeToPB(t *time.Time) *timestamp.Timestamp {
	seconds := t.Unix()
	nanos := int32(t.Sub(time.Unix(seconds, 0)))
	ts := &timestamp.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}

	return ts
}

// NamespacesToPB converts discovery.Namespaces to slices of *discoveryv1.Namespace.
func NamespacesToPB(s discovery.Namespaces) []*discoveryv1.Namespace {
	result := make([]*discoveryv1.Namespace, 0, len(s))

	for i := range s {
		result = append(result, NamespaceToPB(&s[i]))
	}

	return result
}

// NamespacesFromPB converts slices *discoveryv1.Namespace discovery.Namespaces.
func NamespacesFromPB(s []*discoveryv1.Namespace) discovery.Namespaces {
	result := make(discovery.Namespaces, 0, len(s))

	for i := range s {
		result = append(result, *NamespaceFromPB(s[i]))
	}

	return result
}
