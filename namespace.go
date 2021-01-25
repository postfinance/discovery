package discovery

import (
	"errors"
	"sort"
	"time"
)

// ExportConfig defines how services in a namespaces are exported
// for file-based discovery.
type ExportConfig int

// All possible export configurations:
//
// Standard: standard configuration
// Blackbox: for blackbox configurations
// Disabled: no export
const (
	Disabled ExportConfig = iota // disabled
	Standard                     // standard
	Blackbox                     // blackbox
)

// Namespace represents a namespace.
type Namespace struct {
	Name     string       `json:"name"`
	Export   ExportConfig `json:"export"`
	Modified time.Time    `json:"modified,omitempty"`
}

// Validate checks if a services values are valid.
func (n Namespace) Validate() error {
	if n.Name == "" {
		return errors.New("name cannot be empty")
	}

	if !nameRegexp.MatchString(n.Name) {
		return errors.New("name must only contain 'a-z', 'A-Z', '0-9' and '_'")
	}

	return nil
}

// Header creates the header for csv or table output.
func (n Namespace) Header() []string {
	return []string{"NAME", "EXPORTCONFIG", "MODIFIED"}
}

// Row creates a row for csv or table output.
func (n Namespace) Row() []string {
	return []string{n.Name, n.Export.String(), n.Modified.String()}
}

// Namespaces is a list of namespaces.
type Namespaces []Namespace

// Index returns the index of a namespace with name. If not found -1
// is returned.
func (n Namespaces) Index(name string) int {
	for i := range n {
		if name == n[i].Name {
			return i
		}
	}

	return -1
}

// SortByName sorts servers by name.
func (n Namespaces) SortByName() {
	sort.Slice(n, func(i, j int) bool {
		return n[i].Name < n[j].Name
	})
}

// SortByDate sorts servers by registerdate.
func (n Namespaces) SortByDate() {
	sort.Slice(n, func(i, j int) bool {
		return n[j].Modified.Before(n[i].Modified)
	})
}

// Names returns the server names as slice of strings.
func (n Namespaces) Names() []string {
	namespaces := make([]string, 0, len(n))

	for _, namespace := range n {
		namespaces = append(namespaces, namespace.Name)
	}

	return namespaces
}

// DefaultNamespace is used when no namespace is given.
func DefaultNamespace() *Namespace {
	return &Namespace{
		Name:     "default",
		Export:   Standard,
		Modified: time.Now(),
	}
}
