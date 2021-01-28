// Package exporter exports services to filesystem for prometheus file discovery.
package exporter

import (
	"context"
	"os"
	"path/filepath"
	"sync"
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
	serviceRepo  serviceChanLister
	serverRepo   serverGetter
	log          *zap.SugaredLogger
	destDir      string
	server       string
	destinations files
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

	e.server = server

	if err := os.MkdirAll(filepath.Join(e.destDir, server), dirPermissions); err != nil {
		return err
	}

	errHandler := func(err error) {
		e.log.Fatalw("failed to create watcher", "err", err)
	}
	events := e.serviceRepo.Chan(ctx, errHandler)
	ticker := time.NewTicker(reSyncInterval)

	e.log.Info("sync all services")

	if err := e.sync(); err != nil {
		return err
	}

	for {
		select {
		case se, ok := <-events:
			if !ok {
				e.log.Infow("stopping exporter")

				return nil
			}

			e.handle(se)
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

func (e Exporter) handle(event *repo.ServiceEvent) {
	ignore := !event.Service.HasServer(e.server)
	msg := append([]interface{}{
		"event", event.Event.String(),
		"ignore", ignore,
	}, event.Service.KeyVals()...)

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

func (e Exporter) directory() string {
	return filepath.Join(e.destDir, e.server)
}

func (e Exporter) sync() error {
	svcs, err := e.serviceRepo.List("", "")
	if err != nil {
		return err
	}

	// TODO(refacor):
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

type serverGetter interface {
	Get(serverName string) (*discovery.Server, error)
}

type namespaceGetter interface {
	Get(namespaceName string) (*discovery.Namespace, error)
}
