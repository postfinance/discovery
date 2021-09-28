// Package convert converts proto types to domain types and
// vice versa.
package convert

import (
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/exporter"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1"
)

// NamespaceToPB converts *discovery.Namespace to *discoveryv1.Namespace.
func NamespaceToPB(n *discovery.Namespace) *discoveryv1.Namespace {
	pb := &discoveryv1.Namespace{
		Name:     n.Name,
		Export:   int32(n.Export),
		Modified: TimeToPB(&n.Modified),
	}

	return pb
}

// NamespaceFromPB converts *discovery.Namespace to *discoveryv1.Namespace.
func NamespaceFromPB(pb *discoveryv1.Namespace) *discovery.Namespace {
	n := &discovery.Namespace{
		Name:     pb.Name,
		Export:   discovery.ExportConfig(pb.Export),
		Modified: TimeFromPB(pb.Modified),
	}

	return n
}

// ServerToPB converts *discovery.Server to *discoveryv1.Server.
func ServerToPB(s *discovery.Server) *discoveryv1.Server {
	pb := &discoveryv1.Server{
		Name:     s.Name,
		Labels:   s.Labels,
		Modified: TimeToPB(&s.Modified),
		State:    int64(s.State),
	}

	return pb
}

// ServerFromPB converts *discovery.Server to *discoveryv1.Server.
func ServerFromPB(pb *discoveryv1.Server) *discovery.Server {
	s := &discovery.Server{
		Name:     pb.GetName(),
		Labels:   pb.GetLabels(),
		Modified: TimeFromPB(pb.Modified),
		State:    discovery.ServerState(pb.GetState()),
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
		Modified:    TimeToPB(&s.Modified),
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
		Modified:    TimeFromPB(pb.GetModified()),
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

// TimeFromPB converts protobuf timestamps to time.Time.
func TimeFromPB(ts *timestamp.Timestamp) time.Time {
	var t time.Time

	if ts == nil {
		t = time.Unix(0, 0).UTC() // treat nil like the empty Timestamp
	} else {
		t = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	}

	return t
}

// TimeToPB converts time.Time to protobuf timestamps.
func TimeToPB(t *time.Time) *timestamp.Timestamp {
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

// TargetGroupToPB converts *exporter.TargetGroup to discoveryv1.TargetGroup.
func TargetGroupToPB(t *exporter.TargetGroup) *discoveryv1.TargetGroup {
	pb := &discoveryv1.TargetGroup{
		Targets: t.Targets,
		Labels:  t.Labels,
	}

	return pb
}

// TargetGroupFromPB converts *discoveryv1.TargetGroup to *exporter.TargetGroup.
func TargetGroupFromPB(t *discoveryv1.TargetGroup) *exporter.TargetGroup {
	return &exporter.TargetGroup{
		Targets: t.GetTargets(),
		Labels:  t.GetLabels(),
	}
}

// TargetGroupsToPB converts []exporter.TargetGroup to []*discoveryv1.TargetGroup
func TargetGroupsToPB(t []exporter.TargetGroup) []*discoveryv1.TargetGroup {
	result := make([]*discoveryv1.TargetGroup, 0, len(t))

	for i := range t {
		result = append(result, TargetGroupToPB(&t[i]))
	}

	return result
}
