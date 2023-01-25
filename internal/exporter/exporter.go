// Package exporter exports services to filesystem for prometheus file discovery.
package exporter

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/repo"
	"github.com/postfinance/store"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	dirPermissions = 0o750
)

// Exporter writes services from store to filesystem for prometheus
// file based service discovery.
type Exporter struct {
	config               Config
	server               string
	serviceRepo          serviceChanLister
	serverRepo           serverChanGetter
	namespaceRepo        namespaceListGetter
	log                  *zap.SugaredLogger
	httpServer           *http.Server
	serviceWatchEvents   *prometheus.CounterVec
	destinations         files
	serviceWatchDisabled int32
}

// Config configures the exporter.
type Config struct {
	Directory          string
	ResyncInterval     time.Duration
	PrometheusRegistry prometheus.Registerer
	HTTPListenAddr     string
}

// New creates a new exporter.
func New(b store.Backend, log *zap.SugaredLogger, cfg Config) *Exporter {
	namespaceGetter := repo.NewNamespace(b)

	return &Exporter{
		config:        cfg,
		serviceRepo:   repo.NewService(b),
		serverRepo:    repo.NewServer(b),
		namespaceRepo: namespaceGetter,
		serviceWatchEvents: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "discovery_exporter_service_watch_events_total",
			Help: "The total number of service watch events partitioned by operation",
		}, []string{"op"}),
		log: log,
		destinations: files{
			m:               &sync.Mutex{},
			files:           map[string]*file{},
			log:             log,
			namespaceGetter: namespaceGetter,
		},
	}
}

// Start starts the exporter and it blocks until context is done.
func (e *Exporter) Start(ctx context.Context, server string) error {
	errChan := make(chan error)

	if e.config.isStartHTTPServer() {
		go func() {
			if err := e.startHTTP(); err != nil {
				errChan <- err
			}

			close(errChan)
		}()
	}

	_, err := e.serverRepo.Get(server)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return errors.Errorf("server %s not found", server)
		}

		return err
	}

	e.enableWatch()
	e.server = server

	if err := e.createExportDirectories(dirPermissions); err != nil {
		return err
	}

	serviceEvents := e.serviceRepo.Chan(ctx, e.watchErrorHandler)
	serverEvents := e.serverRepo.Chan(ctx, e.watchErrorHandler)
	ticker := time.NewTicker(e.config.ResyncInterval)

	e.log.Info("sync services")

	if err := e.sync(); err != nil {
		return err
	}

	for {
		select {
		case se, ok := <-serverEvents:
			if !ok {
				e.log.Infow("exporter stopped")

				return e.stopHTTP()
			}

			msg := append([]interface{}{"type", se.Event.String()}, se.Server.KeyVals()...)
			e.log.Infow("server event", msg...)

			e.handleServer(se)
		case se, ok := <-serviceEvents:
			if !ok {
				e.log.Infow("stopping exporter")

				return nil
			}

			e.handleService(se)
		case <-ctx.Done():
			e.log.Infow("exporter stopped")

			return e.stopHTTP()
		case <-ticker.C:
			e.log.Debug("initiating resync")
			e.destinations.reset()

			if err := e.sync(); err != nil {
				e.log.Errorw("sync failed", "err", err)
			}
		case err := <-errChan:
			return err
		}
	}
}

func (e *Exporter) handleService(event *repo.ServiceEvent) {
	ignore := !event.Service.HasServer(e.server) && event.Event == repo.Change
	msg := append([]interface{}{
		"event", event.Event.String(),
		"ignore", ignore,
	}, event.Service.KeyVals()...)

	if e.isWatchDisabled() {
		e.log.Debug("watch disabled")
		return
	}

	e.log.Debugw("service event", msg...)

	if ignore {
		return
	}

	if e.serviceWatchEvents != nil {
		e.serviceWatchEvents.WithLabelValues(event.Event.String()).Inc()
	}

	switch event.Event {
	case repo.Change:
		existing := e.destinations.getService(event.Namespace, event.ID)
		// handle path changes
		if existing != nil && existing.Name != event.Service.Name {
			if err := e.destinations.delService(existing.Namespace, existing.ID); err != nil {
				e.log.Errorw("failed to delete service", "namespace", existing.Namespace, "id", existing.ID, "err", err)
			}
		}

		if err := e.destinations.addService(&event.Service); err != nil {
			msg = append(msg, "err", err)
			e.log.Errorw("failed to add service", msg...)
		}
	case repo.Delete:
		// the only possible error is a 'not found' error. we must ingore this error,
		// otherwise we log errors for services that were never registered on this exporter's
		// chache (different server).
		_ = e.destinations.delService(event.Service.Namespace, event.Service.ID)
	default:
		e.log.Errorw("unsupported event", "event", event.Event)
	}

	if err := e.destinations.write(e.directory()); err != nil {
		e.log.Errorw("failed to write discovery files", "err", err)
	}
}

// before deleting or adding a server, the registry disables it. This
// gives us the chance to disable service whatch.
func (e *Exporter) handleServer(event *repo.ServerEvent) {
	switch event.Event {
	case repo.Change:
		if event.Server.State == discovery.Leaving || event.Server.State == discovery.Joining {
			e.log.Debugw("disabling service watcher")
			e.disableWatch()

			return
		}

		if event.Server.State == discovery.Active {
			e.log.Debug("sync services")

			e.destinations.reset()

			if err := e.sync(); err != nil {
				e.log.Errorw("sync failed", "err", err)
			}

			e.log.Debugw("enabling service watcher")
			e.enableWatch()
		}
	case repo.Delete:
		e.log.Debug("sync services")

		e.destinations.reset()

		if err := e.sync(); err != nil {
			e.log.Errorw("sync failed", "err", err)
		}

		e.log.Debugw("enabling service watcher")
		e.enableWatch()
	default:
		e.log.Errorw("unsupported event", "event", event.Event)
	}
}

func (e *Exporter) directory() string {
	return filepath.Join(e.config.Directory, e.server)
}

func (e *Exporter) isWatchDisabled() bool {
	return atomic.LoadInt32(&e.serviceWatchDisabled) == 1
}

func (e *Exporter) disableWatch() {
	atomic.StoreInt32(&e.serviceWatchDisabled, 1)
}

func (e *Exporter) enableWatch() {
	atomic.StoreInt32(&e.serviceWatchDisabled, 0)
}

func (e *Exporter) createExportDirectories(permission fs.FileMode) error {
	namespaces, err := e.namespaceRepo.List()
	if err != nil {
		return err
	}

	for _, n := range namespaces {
		dir := filepath.Join(e.config.Directory, e.server, n.Name)
		e.log.Infow("creating export directory", "path", dir)

		if err := os.MkdirAll(dir, permission); err != nil {
			return err
		}
	}

	return nil
}

func (e *Exporter) watchErrorHandler(err error) {
	e.log.Fatalw("failed to create watcher", "err", err)
}

func (e *Exporter) sync() error {
	svcs, err := e.serviceRepo.List("", "")
	if err != nil {
		return err
	}

	for i := range svcs {
		if !svcs[i].HasServer(e.server) {
			continue
		}

		if err := e.destinations.addService(&svcs[i]); err != nil {
			return err
		}
	}

	if err := e.destinations.delObsoleteFiles(e.directory()); err != nil {
		return err
	}

	return e.destinations.write(e.directory())
}

type serviceChanLister interface {
	List(namespace, selector string) (discovery.Services, error)
	Chan(context.Context, func(error)) <-chan *repo.ServiceEvent
}

type serverChanGetter interface {
	Get(serverName string) (*discovery.Server, error)
	Chan(context.Context, func(error)) <-chan *repo.ServerEvent
}

type namespaceListGetter interface {
	namespaceGetter
	List() (discovery.Namespaces, error)
}

type namespaceGetter interface {
	Get(namespaceName string) (*discovery.Namespace, error)
}

func (cfg Config) isStartHTTPServer() bool {
	return cfg.PrometheusRegistry != nil && cfg.HTTPListenAddr != ""
}
