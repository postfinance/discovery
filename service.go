// Package discovery contains domain logic for service discovery.
package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/labels"
)

var (
	nameRegexp = regexp.MustCompile("^[[:alnum:]_-]+$")
)

// Service contains all information for service discovery.
type Service struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Namespace   string    `json:"namespace,omitempty"`
	Endpoint    *url.URL  `json:"endpoint,omitempty"`
	Selector    string    `json:"selector,omitempty"`
	Servers     []string  `json:"servers,omitempty"`
	Labels      Labels    `json:"labels,omitempty"`
	Description string    `json:"description,omitempty"`
	Modified    time.Time `json:"modified,omitempty"`
}

// NewService creates a new service with ID and timestamp.
func NewService(name, endpoint string) (*Service, error) {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, fmt.Errorf("'%s' is not a valid url: %v", endpoint, err)
	}

	s := &Service{
		Name:      name,
		Namespace: DefaultNamespace().Name,
		Endpoint:  u,
		Modified:  time.Now(),
	}

	if err := s.Validate(); err != nil {
		return nil, err
	}

	return s, nil
}

// MustNewService panics if endpoint or name is not valid.
func MustNewService(name, endpoint string) *Service {
	s, err := NewService(name, endpoint)
	if err != nil {
		panic(err)
	}

	return s
}

// Validate checks if a services values are valid.
func (s Service) Validate() error {
	if s.Name == "" {
		return errors.New("name cannot be empty")
	}

	if !nameRegexp.MatchString(s.Name) {
		return errors.New("name must only contain 'a-z', 'A-Z', '0-9', '-' and '_'")
	}

	if s.Namespace == "" {
		return errors.New("namespace cannot be empty")
	}

	if !nameRegexp.MatchString(s.Namespace) {
		return errors.New("namespace must only contain 'a-z', 'A-Z', '0-9', '-' and '_'")
	}

	if s.Endpoint == nil {
		return errors.New("endpoint cannot be null")
	}

	return nil
}

// HasServer returns true if service has serverName in its Servers slice.
func (s Service) HasServer(serverName string) bool {
	for _, name := range s.Servers {
		if name == serverName {
			return true
		}
	}

	return false
}

// Header creates the header for csv or table output.
func (s Service) Header() []string {
	return []string{"NAME", "NAMESPACE", "ID", "ENDPOINT", "SERVERS", "LABELS", "SELECTOR", "MODIFIED", "DESCRIPTION"}
}

// Row creates a row for csv or table output.
func (s Service) Row() []string {
	return []string{s.Name, s.Namespace, s.ID, s.Endpoint.String(), strings.Join(s.Servers, ","), s.Labels.String(), s.Selector, s.Modified.Format(time.RFC3339), s.Description}
}

// UnmarshalJSON is a custom json unmarshaller.
func (s *Service) UnmarshalJSON(j []byte) error {
	raw := struct {
		ID          string    `json:"id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Namespace   string    `json:"namespace,omitempty"`
		Endpoint    string    `json:"endpoint,omitempty"`
		Labels      Labels    `json:"labels,omitempty"`
		Servers     []string  `json:"servers,omitempty"`
		Selector    string    `json:"selector,omitempty"`
		Description string    `json:"description,omitempty"`
		Modified    time.Time `json:"modified,omitempty"`
	}{}

	err := json.Unmarshal(j, &raw)
	if err != nil {
		return err
	}

	s.ID = raw.ID
	s.Name = raw.Name
	s.Namespace = raw.Namespace
	s.Labels = raw.Labels
	s.Selector = raw.Selector
	s.Servers = raw.Servers
	s.Description = raw.Description
	s.Modified = raw.Modified

	if raw.Endpoint == "" {
		s.Endpoint = nil
		return nil
	}

	u, err := url.Parse(raw.Endpoint)
	if err != nil {
		return err
	}

	s.Endpoint = u

	return nil
}

// MarshalJSON is a custom JSON marshaler.
func (s Service) MarshalJSON() ([]byte, error) {
	ep := ""
	if s.Endpoint != nil {
		ep = s.Endpoint.String()
	}

	raw := struct {
		ID          string    `json:"id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Namespace   string    `json:"namespace,omitempty"`
		Endpoint    string    `json:"endpoint,omitempty"`
		Labels      Labels    `json:"labels,omitempty"`
		Servers     []string  `json:"servers,omitempty"`
		Selector    string    `json:"selector,omitempty"`
		Description string    `json:"description,omitempty"`
		Modified    time.Time `json:"modified,omitempty"`
	}{
		ID:          s.ID,
		Name:        s.Name,
		Namespace:   s.Namespace,
		Endpoint:    ep,
		Labels:      s.Labels,
		Servers:     s.Servers,
		Selector:    s.Selector,
		Description: s.Description,
		Modified:    s.Modified,
	}

	return json.Marshal(raw)
}

// Services is a slice of Services
type Services []Service

// FilterFunc is a function to filter services. If function returns true
// service is selected else omitted.
type FilterFunc func(Service) bool

// Filter filters Services with FilterFunc.
func (s Services) Filter(filters ...FilterFunc) Services {
	services := Services{}

	for i := range s {
		selectService := true
		for _, f := range filters {
			selectService = selectService && f(s[i])
		}

		if selectService {
			services = append(services, s[i])
		}
	}

	return services
}

// ServiceByName filters Services by Name.
func ServiceByName(r *regexp.Regexp) FilterFunc {
	return func(s Service) bool {
		return r.MatchString(s.Name)
	}
}

// ServiceByEndpoint filters Services by Endpoint.
func ServiceByEndpoint(r *regexp.Regexp) FilterFunc {
	return func(s Service) bool {
		return r.MatchString(s.Endpoint.String())
	}
}

// ServiceByServer filters Services by Server.
func ServiceByServer(r *regexp.Regexp) FilterFunc {
	return func(s Service) bool {
		for _, srv := range s.Servers {
			if r.MatchString(srv) {
				return true
			}
		}

		return false
	}
}

// ServiceBySelector filters Services by Selector.
func ServiceBySelector(selector labels.Selector) FilterFunc {
	return func(s Service) bool {
		return selector.Matches(s.Labels)
	}
}

// SortByEndpoint sorts services by endpoint.
func (s Services) SortByEndpoint() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Endpoint.String() < s[j].Endpoint.String()
	})
}

// SortByDate sorts servers by modification date.
func (s Services) SortByDate() {
	sort.Slice(s, func(i, j int) bool {
		return s[j].Modified.Before(s[i].Modified)
	})
}

// Labels represents key value pairs.
type Labels map[string]string

func (l Labels) String() string {
	keys := []string{}

	for k := range l {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	s := ""
	sep := ""

	for _, k := range keys {
		if len(s) > 0 {
			sep = ","
		}

		s = s + sep + k + "=" + l[k]
	}

	return s
}

// Has returs true if Labels contains key.
func (l Labels) Has(key string) bool {
	_, ok := l[key]
	return ok
}

// Get gets the value for key.
func (l Labels) Get(key string) string {
	return l[key]
}

// KeyVals represents the service as slice of interface.
func (s Service) KeyVals() []interface{} {
	return []interface{}{
		"id", s.ID,
		"name", s.Name,
		"namespace", s.Namespace,
		"endpoint", s.Endpoint,
		"modified", s.Modified,
		"selector", s.Selector,
		"description", s.Description,
	}
}
