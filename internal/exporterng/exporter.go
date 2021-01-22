// +build bla

// Package exporter exports services to filesystem for prometheus file discovery.
package exporter

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/renameio"
	"github.com/postfinance/store"
	"github.com/zbindenren/discovery"
	"github.com/zbindenren/discovery/internal/repo"
	"go.uber.org/zap"
)

const (
	jobLabel         = "job"
	instanceLabel    = "instance"
	schemeLabel      = "__scheme__"
	metricsPathLabel = "__metrics_path__"
)

// Exporter writes services from store to filesystem for prometheus
// file based service discovery.
type Exporter struct {
	repo         serviceChanLister
	log          *zap.SugaredLogger
	directory    string
	destinations files
}

// New creates a new exporter.
func New(b store.Backend, log *zap.SugaredLogger, directory string) *Exporter {
	return &Exporter{
		repo:         repo.NewService(b),
		log:          log,
		directory:    directory,
		destinations: files{},
	}
}

// Start starts the exporter and it blocks until context is done.
func (e Exporter) Start(ctx context.Context, reSyncInterval time.Duration) error {
	errHandler := func(err error) {
		e.log.Fatalw("failed to create watcher", "err", err)
	}

	if err := e.sync(); err != nil {
		return err
	}

	events := e.repo.Chan(ctx, errHandler)
	ticker := time.NewTicker(reSyncInterval)

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
			if err := e.sync(); err != nil {
				e.log.Errorw("sync failed", "err", err)
			}
		}
	}
}

func (e *Exporter) handle(se *repo.ServiceEvent) {
	msg := append([]interface{}{
		"event", se.Event.String(),
	}, se.Service.KeyVals()...)

	e.log.Infow("service event", msg...)

	switch se.Event {
	case repo.Change:
		tgs := newTargetGroups(e.directory, &se.Service)
		if len(tgs) == 0 {
			e.log.Errorw("no servers registered")
		}

		e.destinations.add(e.log, tgs...)
	case repo.Delete:
		// on delete only have namespace and id information is available
		tgs := e.destinations.findAllTargetGroup(se.Service.Namespace, se.Service.ID)
		for _, tg := range tgs {
			e.destinations.rm(e.log, tg.server, tg.namespace, tg.id)
		}
	default:
		e.log.Error("unsupported event", "event", se.Event)
	}

	e.destinations.write(e.log)
}

func (e *Exporter) sync() error {
	e.log.Infow("sync all services")

	svcs, err := e.repo.List("", "")
	if err != nil {
		return err
	}

	for i := range svcs {
		s := svcs[i]
		for _, tg := range newTargetGroups(e.directory, &s) {
			e.destinations.add(e.log, tg)
		}
	}

	syncedFiles := []string{}
	err = filepath.Walk(e.directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			d, err := ioutil.ReadFile(path) // nolint: gosec
			if err != nil {
				return err
			}

			if _, ok := e.destinations[path]; ok {
				e.destinations[path].sum = fmt.Sprintf("%x", sha256.Sum256(d))
			}

			syncedFiles = append(syncedFiles, path)

			return nil
		})

	if err != nil {
		return err
	}

	for _, f := range syncedFiles {
		if _, ok := e.destinations[f]; !ok {
			e.log.Infow("removing file", "path", f)

			if err := os.Remove(f); err != nil {
				return err
			}
		}
	}

	e.destinations.write(e.log)

	return nil
}

// TargetGroup represents a prometheus target group for file service discovery.
type targetGroup struct {
	Targets   []string         `json:"targets,omitempty"`
	Labels    discovery.Labels `json:"labels,omitempty"`
	id        string
	namespace string
	server    string
	name      string
	directory string
}

func (tg targetGroup) path() string {
	return createAbsPath(tg.directory, tg.server, tg.namespace, tg.name)
}

func createAbsPath(dstDir, server, namespace, jobName string) string {
	return path.Join(dstDir, server, namespace, jobName+".json")
}

type files map[string]*file

func (f files) add(l *zap.SugaredLogger, tgs ...targetGroup) {
	for _, tg := range tgs {
		// if service job changes the path changes and we have
		// to delete targetGroup from old path.
		existing := f.findTargetGroup(tg.server, tg.namespace, tg.id)
		// handle path changes
		if existing != nil && existing.name != tg.name {
			f.rm(l, existing.server, existing.namespace, existing.id)
		}

		if _, ok := f[tg.path()]; !ok {
			f[tg.path()] = &file{
				path:         tg.path(),
				targetGroups: targetGroups{},
			}
		}

		f[tg.path()].targetGroups[tg.id] = tg
	}
}

func (f files) rm(l *zap.SugaredLogger, server, namespace, id string) {
	tg := f.findTargetGroup(server, namespace, id)
	if tg == nil {
		l.Errorw("no target group found",
			"server", server,
			"namespace", namespace,
			"id", id,
		)

		return
	}

	delete(f[tg.path()].targetGroups, tg.id)
	f[tg.path()].sum = ""

	if len(f[tg.path()].targetGroups) == 0 {
		l.Infow("removing discovery file", "path", tg.path())

		if err := os.Remove(tg.path()); err != nil {
			l.Errorw("failed to remove file", "path", tg.path(), "err", err)
		}
	}
}

func (f files) findAllTargetGroup(namespace, id string) []targetGroup {
	result := []targetGroup{}

	for absPath := range f {
		if tg, ok := f[absPath].targetGroups[id]; ok {
			if tg.namespace == namespace {
				result = append(result, tg)
			}
		}
	}

	return result
}

func (f files) findTargetGroup(server, namespace, id string) *targetGroup {
	for absPath := range f {
		if tg, ok := f[absPath].targetGroups[id]; ok {
			if tg.namespace == namespace && tg.server == server {
				return &tg
			}
		}
	}

	return nil
}

func (f files) write(l *zap.SugaredLogger) {
	for _, file := range f {
		if err := file.write(l); err != nil {
			l.Errorw("failed to write file", "err", err)
		}
	}
}

type file struct {
	path         string
	sum          string
	targetGroups targetGroups
}

func (f *file) write(l *zap.SugaredLogger) error {
	if len(f.targetGroups) == 0 {
		return nil
	}

	data, err := json.Marshal(f.targetGroups.list())
	if err != nil {
		return err
	}

	sum := fmt.Sprintf("%x", sha256.Sum256(data))
	if sum == f.sum {
		l.Debugw("discovery file already up-to-date", "path", f.path)
		return nil
	}

	dir := filepath.Dir(f.path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	l.Infow("updating discovery file", "path", f.path, "sum", sum, "old", f.sum)

	f.sum = sum

	return renameio.WriteFile(f.path, data, 0644) // nolint: gocritic
}

func newTargetGroups(directory string, s *discovery.Service) []targetGroup {
	tg := targetGroup{
		id:        s.ID,
		namespace: s.Namespace,
		name:      s.Name,
		directory: directory,
		Targets:   []string{s.Endpoint.Host},
		Labels: discovery.Labels{
			jobLabel:         s.Name,
			instanceLabel:    s.Endpoint.Host,
			schemeLabel:      s.Endpoint.Scheme,
			metricsPathLabel: strings.TrimRight(s.Endpoint.Path, "/"),
		},
	}

	for k, v := range s.Labels {
		tg.Labels[k] = v
	}

	for k, v := range s.Endpoint.Query() {
		tg.Labels["__param_"+k] = v[0]
	}

	result := make([]targetGroup, 0, len(s.Servers))

	for _, server := range s.Servers {
		tg.server = server
		result = append(result, tg)
	}

	return result
}

type targetGroups map[string]targetGroup

func (t targetGroups) list() []targetGroup {
	l := make([]targetGroup, 0, len(t))

	for _, v := range t {
		l = append(l, v)
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].Targets[0] < l[j].Targets[0]
	})

	return l
}

type serviceChanLister interface {
	List(namespace, selector string) (discovery.Services, error)
	Chan(context.Context, func(error)) <-chan *repo.ServiceEvent
}
