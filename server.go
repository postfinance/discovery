package discovery

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/labels"
)

// Server represents a registered server.
//
// With kubernetes selectors it is possible to select a server by labels.
// If IsActive is false, no services are distributed to this server.
type Server struct {
	Name      string    `json:"name"`
	Labels    Labels    `json:"labels"`
	IsEnabled bool      `json:"enabled"`
	Modified  time.Time `json:"modified,omitempty"`
}

// NewServer creates a new server instance.
func NewServer(name string, l Labels) *Server {
	return &Server{
		Name:      name,
		Labels:    l,
		Modified:  time.Now(),
		IsEnabled: true,
	}
}

// Validate checks if a services values are valid.
func (s Server) Validate() error {
	if s.Name == "" {
		return errors.New("name cannot be empty")
	}

	return nil
}

// Servers is a list of servers.
type Servers []Server

// SortByName sorts servers by name.
func (s Servers) SortByName() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Name < s[j].Name
	})
}

// SortByDate sorts servers by registerdate.
func (s Servers) SortByDate() {
	sort.Slice(s, func(i, j int) bool {
		return s[j].Modified.Before(s[i].Modified)
	})
}

// Enabled filters out disabled servers.
func (s Servers) Enabled() Servers {
	servers := make(Servers, 0, len(s))

	for _, server := range s {
		if server.IsEnabled {
			servers = append(servers, server)
		}
	}

	return servers
}

// Names returns the server names as slice of strings.
func (s Servers) Names() []string {
	servers := make([]string, 0, len(s))

	for _, server := range s {
		servers = append(servers, server.Name)
	}

	return servers
}

// Filter filters Servers.
func (s Servers) Filter(f func(Server) bool) Servers {
	servers := Servers{}

	for _, server := range s {
		if f(server) {
			servers = append(servers, server)
		}
	}

	return servers
}

// ServersBySelector filters Servers by Selector.
func ServersBySelector(selector labels.Selector) func(Server) bool {
	return func(s Server) bool {
		return selector.Matches(s.Labels)
	}
}

// Header creates the header for csv or table output.
func (s Server) Header() []string {
	return []string{"NAME", "MODIFIED", "ENABLED", "LABELS"}
}

// Row creates a row for csv or table output.
func (s Server) Row() []string {
	return []string{s.Name, s.Modified.Format(time.RFC3339), fmt.Sprintf("%v", s.IsEnabled), s.Labels.String()}
}
