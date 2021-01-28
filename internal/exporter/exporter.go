// Package exporter exports services to filesystem for prometheus file discovery.
package exporter

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/postfinance/store"
	"github.com/zbindenren/discovery"
	"github.com/zbindenren/discovery/internal/repo"
	"go.uber.org/zap"
)

const (
	dirPermissions = 0750
)

// Exporter writes services from store to filesystem for prometheus
// file based service discovery.
type Exporter struct {
	serviceRepo          serviceChanLister
	serverRepo           serverChanGetter
	log                  *zap.SugaredLogger
	destDir              string
	server               string
	destinations         files
	serviceWatchDisabled int32
}

// New creates a new exporter.
func New(b store.Backend, log *zap.SugaredLogger, destDir string) *Exporter {
	return &Exporter{
		serviceRepo: repo.NewService(b),
		serverRepo:  repo.NewServer(b),
		log:         log,
		destDir:     destDir,
		destinations: files{
			m:               &sync.Mutex{},
			files:           map[string]*file{},
			log:             log,
			namespaceGetter: repo.NewNamespace(b),
		},
	}
}

// Start starts the exporter and it blocks until context is done.
func (e *Exporter) Start(ctx context.Context, server string, reSyncInterval time.Duration) error {
	_, err := e.serverRepo.Get(server)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return errors.Errorf("server %s not found", server)
		}

		return err
	}

	e.enableWatch()
	e.server = server

	if err := os.MkdirAll(filepath.Join(e.destDir, server), dirPermissions); err != nil {
		return err
	}

	errHandler := func(err error) {
		e.log.Fatalw("failed to create watcher", "err", err)
	}
	serviceEvents := e.serviceRepo.Chan(ctx, errHandler)
	serverEvents := e.serverRepo.Chan(ctx, errHandler)
	ticker := time.NewTicker(reSyncInterval)

	e.log.Info("sync services")

	if err := e.sync(); err != nil {
		return err
	}

	for {
		select {
		case se, ok := <-serverEvents:
			if !ok {
				e.log.Infow("stopping exporter")

				return nil
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
			e.log.Infow("stopping exporter")

			return nil
		case <-ticker.C:
			e.log.Debug("initiating resync")
			e.destinations.reset()

			if err := e.sync(); err != nil {
				e.log.Errorw("sync failed", "err", err)
			}
		}
	}
}

func (e *Exporter) handleService(event *repo.ServiceEvent) {
	ignore := !event.Service.HasServer(e.server)
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
		if err := e.destinations.delService(event.Service.Namespace, event.Service.ID); err != nil {
			e.log.Errorw("failed to delete service", "namespace", event.Service.Namespace, "id", event.Service.ID, "err", err)
		}
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
	return filepath.Join(e.destDir, e.server)
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

type namespaceGetter interface {
	Get(namespaceName string) (*discovery.Namespace, error)
}
