// Package registry is responsible for registering and unregistering
// services and servers.
package registry

import (
	"context"
	"errors"
	"fmt"
	"hash/crc64"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/hash"
	"github.com/postfinance/discovery/internal/repo"
	"github.com/postfinance/store"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// Common errors
var (
	ErrNotFound         = errors.New("not found")
	ErrContainsServices = errors.New("server has registered services")
)

// Registry registers server or service.
type Registry struct {
	log            *zap.SugaredLogger
	serverRepo     *repo.Server
	serviceRepo    *repo.Service
	namespaceRepo  *repo.Namespace
	jumpHasher     *hash.Jump
	idGenerator    func(string) string
	numReplicas    int
	servicesCount  *prometheus.GaugeVec
	namespaceCache namespaceCache
}

// New creates a new registry.
func New(backend store.Backend, reg prometheus.Registerer, log *zap.SugaredLogger, numReplicas int) (*Registry, error) {
	if numReplicas < 1 {
		return nil, errors.New("number of replicas has to be >= 1")
	}

	servicesCount := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "discovery_services_count",
			Help: "Number of registered services per server and namespace.",
		},
		[]string{"server", "namespace"},
	)

	reg.MustRegister(servicesCount)

	registry := Registry{
		log:           log,
		jumpHasher:    hash.New(crc64.New(crc64.MakeTable(0xC96C5795D7870F42))),
		idGenerator:   repo.IDGenerator(),
		numReplicas:   numReplicas,
		serverRepo:    repo.NewServer(backend),
		serviceRepo:   repo.NewService(backend),
		namespaceRepo: repo.NewNamespace(backend),
		servicesCount: servicesCount,
		namespaceCache: namespaceCache{
			m:          &sync.Mutex{},
			namespaces: map[string]discovery.Namespace{},
		},
	}

	if err := registry.initNamespaceCache(); err != nil {
		return nil, err
	}

	return &registry, nil
}

// StartCacheUpdater starts a namespace cache updater. It resyncs cache all reSyncInterval.
func (r Registry) StartCacheUpdater(ctx context.Context, reSyncInterval time.Duration) {
	namespaceEventChan := r.namespaceRepo.Chan(ctx, func(err error) {
		r.log.Fatalw("failed to create namespace watcher", "err", err)
	})

	ticker := time.NewTicker(reSyncInterval)

	for {
		select {
		case n := <-namespaceEventChan:
			r.log.Debug("updating namespace cache")

			switch n.Event {
			case repo.Change:
				r.namespaceCache.add(n.Namespace)
			case repo.Delete:
				r.namespaceCache.del(n.Name)
			default:
				r.log.Errorw("unsupported namespace envent type", "event", n.Event.String())
			}
		case <-ctx.Done():
			r.log.Infow("stopping namespace cache updater")

			return
		case <-ticker.C:
			r.log.Debug("initiating namespace cache sync")

			if err := r.initNamespaceCache(); err != nil {
				r.log.Errorw("failed to sync namespace cache", "err", err)
			}
		}
	}
}

// StartServiceCounterUpdater updates service counter metrics every interval. It runs until context ctx
// is canceled.
func (r Registry) StartServiceCounterUpdater(ctx context.Context, interval time.Duration) {
	r.log.Infow("initializing service counter", "interval", interval)

	if err := r.updateServiceCounter(); err != nil {
		r.log.Errorw("failed to initialize service counter", "err", err)
	}

	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			r.log.Info("stopping service counter")

			return
		case <-ticker.C:
			r.log.Debug("updating service counter")

			if err := r.updateServiceCounter(); err != nil {
				r.log.Errorw("service counter update failed", "err", err)
			}
		}
	}
}

// RegisterServer registers a server.
func (r *Registry) RegisterServer(name string, labels discovery.Labels) (*discovery.Server, error) {
	s := discovery.NewServer(name, labels)

	if err := s.Validate(); err != nil {
		return nil, err
	}

	_, err := r.serverRepo.Save(*s)
	if err != nil {
		return nil, fmt.Errorf("failed to save server %s: %w", s.Name, err)
	}

	if _, err := r.ReRegisterAllServices(); err != nil {
		return nil, fmt.Errorf("failed to reregister all services: %w", err)
	}

	s.State = discovery.Active
	r.log.Infow("register server", s.KeyVals()...)

	ns, err := r.serverRepo.Save(*s)
	if err != nil {
		return nil, fmt.Errorf("failed to save server %s: %w", s.Name, err)
	}

	return ns, nil
}

// UnRegisterServer unregisters a server.
func (r *Registry) UnRegisterServer(name string) error {
	s, err := r.serverRepo.Get(name)
	if err != nil {
		return err
	}

	s.State = discovery.Leaving

	if _, err := r.serverRepo.Save(*s); err != nil {
		return err
	}

	r.log.Infow("unregister server", "name", name)

	if _, err := r.ReRegisterAllServices(); err != nil {
		return fmt.Errorf("failed to reregister all services: %w", err)
	}

	if err := r.serverRepo.Delete(name); err != nil {
		return fmt.Errorf("failed to delete server %s: %w", name, err)
	}

	return nil
}

// ListServer lists servers by selector.
func (r *Registry) ListServer(selector string) (discovery.Servers, error) {
	return r.serverRepo.List(selector)
}

// RegisterService registers a service.
func (r *Registry) RegisterService(s discovery.Service) (*discovery.Service, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	if !r.namespaceCache.hasNamespace(s.Namespace) {
		return nil, fmt.Errorf("namespace %s: %w", s.Namespace, ErrNotFound)
	}

	r.log.Infow("register service", s.KeyVals()...)

	servers, err := r.get(s.Endpoint.String(), r.numReplicas, s.Selector)
	if err != nil {
		return nil, err
	}

	if len(servers) == 0 {
		return nil, ErrNotFound
	}

	s.Servers = servers.Names()

	return r.serviceRepo.Save(s)
}

// UnRegisterService removes a service by id or endpoint. If namespace is empty string
// then discovery.DefaultNamespace is used.
func (r *Registry) UnRegisterService(idOrEndpoint, namespace string) error {
	if namespace == "" {
		namespace = discovery.DefaultNamespace().Name
	}

	id := idOrEndpoint

	if strings.Contains(id, ":") {
		id = r.idGenerator(idOrEndpoint)
	}

	r.log.Infow("unregister service", "id", id, "namespace", namespace)

	if err := r.serviceRepo.Delete(id, namespace); err != nil {
		return err
	}

	return nil
}

// ReRegisterAllServices reregisters all services.
func (r *Registry) ReRegisterAllServices() (numChanges int, err error) {
	allServices, err := r.serviceRepo.List("", "")
	if err != nil {
		return 0, err
	}

	for i := range allServices {
		s := allServices[i]

		ns, err := r.RegisterService(s)
		if err != nil {
			return 0, err
		}

		if !reflect.DeepEqual(s.Servers, ns.Servers) {
			numChanges++
		}
	}

	return numChanges, nil
}

// ListService lists all services.
func (r *Registry) ListService(namespace, selector string) (discovery.Services, error) {
	return r.serviceRepo.List(namespace, selector)
}

// RegisterNamespace registers a namespace.
func (r *Registry) RegisterNamespace(n discovery.Namespace) (*discovery.Namespace, error) {
	if err := n.Validate(); err != nil {
		return nil, err
	}

	r.log.Infow("register namespace", "name", n.Name, "exportconfig", n.Export.String())

	n.Modified = time.Now()

	ns, err := r.namespaceRepo.Save(n)
	if err != nil {
		return ns, fmt.Errorf("failed to save namespace %s: %w", n.Name, err)
	}

	r.namespaceCache.add(n)

	return ns, nil
}

// ListNamespaces lists all services.
func (r *Registry) ListNamespaces() (discovery.Namespaces, error) {
	return r.namespaceCache.list(), nil
}

// UnRegisterNamespace unregisters a namespace.
func (r *Registry) UnRegisterNamespace(name string) error {
	r.log.Infow("unregister namespace", "name", name)

	if err := r.serviceRepo.DeleteFromNamespace(name); err != nil {
		return fmt.Errorf("failed to delete all services in namespace %s: %w", name, err)
	}

	if err := r.namespaceRepo.Delete(name); err != nil {
		return fmt.Errorf("failed to delete namespace %s: %w", name, err)
	}

	r.namespaceCache.del(name)

	return nil
}

func (r *Registry) updateServiceCounter() error {
	services, err := r.serviceRepo.List("", "")
	if err != nil {
		return err
	}

	r.servicesCount.Reset()

	for i := range services {
		for _, server := range services[i].Servers {
			r.servicesCount.WithLabelValues(server, services[i].Namespace).Inc()
		}
	}

	return nil
}

// get gets one or numReplica server for a key via consistent hasher. If numReplica is larger
// than the number of servers, len(servers) is used.
func (r *Registry) get(key string, numReplica int, selector string) (discovery.Servers, error) {
	candidates, err := r.serverRepo.List(selector)
	if err != nil {
		return nil, err
	}

	candidates = candidates.Enabled()
	candidates.SortByName()

	if numReplica > len(candidates) {
		numReplica = len(candidates)
	}

	if numReplica == len(candidates) {
		return candidates, nil
	}

	result := make(discovery.Servers, 0, numReplica)
	i := r.jumpHasher.HashString(key, len(candidates))

	for {
		if len(result) == numReplica {
			break
		}

		if i >= len(candidates) {
			i = 0
		}

		result = append(result, candidates[i])

		i++
	}

	result.SortByName()

	return result, nil
}

func (r *Registry) initNamespaceCache() error {
	namespaces, err := r.namespaceRepo.List()
	if err != nil {
		return err
	}

	r.namespaceCache.reset()

	for _, n := range namespaces {
		r.namespaceCache.add(n)
	}

	return nil
}

type namespaceCache struct {
	m          *sync.Mutex
	namespaces map[string]discovery.Namespace
}

func (n *namespaceCache) hasNamespace(name string) bool {
	n.m.Lock()
	_, ok := n.namespaces[name]
	n.m.Unlock()

	return ok
}

func (n *namespaceCache) reset() {
	n.m.Lock()
	n.namespaces = map[string]discovery.Namespace{}
	n.m.Unlock()
}

func (n *namespaceCache) list() discovery.Namespaces {
	n.m.Lock()
	defer n.m.Unlock()

	namespaces := make(discovery.Namespaces, 0, len(n.namespaces))

	for _, ns := range n.namespaces {
		namespaces = append(namespaces, ns)
	}

	namespaces.SortByName()

	return namespaces
}

func (n *namespaceCache) add(namespace discovery.Namespace) {
	n.m.Lock()
	n.namespaces[namespace.Name] = namespace
	n.m.Unlock()
}

func (n *namespaceCache) del(name string) {
	n.m.Lock()
	delete(n.namespaces, name)
	n.m.Unlock()
}
