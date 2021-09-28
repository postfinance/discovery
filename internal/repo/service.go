package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/postfinance/discovery"
	"github.com/postfinance/store"
	uuid "github.com/satori/go.uuid"
	"k8s.io/apimachinery/pkg/labels"
)

// Service represents the service repository.
type Service struct {
	backend store.Backend
	prefix  string
	idGen   func(string) string
	w       *serviceWatcher
}

// NewService creates a new service repo.
func NewService(backend store.Backend) *Service {
	return &Service{
		backend: backend,
		prefix:  servicePrefix,
		idGen:   IDGenerator(),
	}
}

// IDGenerator returns a uuid generator.
func IDGenerator() func(string) string {
	uuidNamespace := uuid.NewV5(uuid.NamespaceDNS, "discovery")

	return func(s string) string {
		return uuid.NewV5(uuidNamespace, s).String()
	}
}

// Get gets a service by id and namespace.
func (s *Service) Get(id, namespace string) (*discovery.Service, error) {
	var svc discovery.Service

	_, err := s.backend.Get(s.key(namespace, id), store.WithHandler(func(k, v []byte) error {
		err := json.Unmarshal(v, &svc)
		if err != nil {
			return err
		}

		return nil
	}))

	if err != nil {
		return nil, err
	}

	return &svc, nil
}

// Save creates or updates a service. It returns the service with the generated id.
func (s *Service) Save(svc discovery.Service) (*discovery.Service, error) {
	newID := s.idGen(svc.Endpoint.String())

	// check if endpoint got changed
	if svc.ID != "" && newID != svc.ID {
		if err := s.Delete(svc.ID, svc.Namespace); err != nil {
			return nil, err
		}
	}

	svc.ID = newID
	svc.Modified = time.Now()

	if _, err := store.Put(s.backend, s.key(svc.Namespace, svc.ID), svc); err != nil {
		return nil, err
	}

	return &svc, nil
}

// Delete removes a service from repo.
func (s *Service) Delete(id, namespace string) error {
	count, err := s.backend.Del(s.key(namespace, id))
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("%s/%s: %w", namespace, id, ErrNotFound)
	}

	return nil
}

// List lists all services for namespace and selector. If namespace is empty string, all
// selected services are returned.
func (s *Service) List(namespace, selector string) (discovery.Services, error) {
	services := discovery.Services{}

	_, err := s.backend.Get(path.Join(s.prefix, namespace), store.WithPrefix(), store.WithHandler(func(k, v []byte) error {
		svc := discovery.Service{}

		err := json.Unmarshal(v, &svc)
		if err != nil {
			return err
		}
		services = append(services, svc)
		return nil
	}))

	if err != nil {
		return nil, err
	}

	if selector == "" {
		return services, nil
	}

	sel, err := labels.Parse(selector)
	if err != nil {
		return nil, err
	}

	result := services.Filter(discovery.ServiceBySelector(sel))

	return result, nil
}

// DeleteFromNamespace deletes all services in namespace.
func (s *Service) DeleteFromNamespace(namespace string) error {
	if namespace == "" {
		return errors.New("namespace cannot be empty string")
	}

	_, err := s.backend.Del(path.Join(s.prefix, namespace), store.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

func serviceFromKey(key []byte) (*discovery.Service, error) {
	k := string(key)
	fields := strings.Split(k, "/")

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid service path: %s", k)
	}

	id := fields[len(fields)-1]
	ns := fields[len(fields)-2]

	return &discovery.Service{
		ID:        id,
		Namespace: ns,
	}, nil
}

func (s *Service) key(namespace, id string) string {
	return path.Join(s.prefix, namespace, id)
}
