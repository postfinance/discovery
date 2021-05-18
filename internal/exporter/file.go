package exporter

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/google/renameio"
	"github.com/postfinance/discovery"
	"go.uber.org/zap"
)

// file represents a prometheus discovery file. It contains
// all target groups for a job.
type file struct {
	job       string
	namespace string
	hash      string
	exportCfg discovery.ExportConfig
	m         *sync.Mutex
	services  services
}

func (f *file) relativePath() string {
	return filepath.Join(f.namespace, f.job) + ".json"
}

func (f *file) addService(s service) error {
	if err := f.checkService(s); err != nil {
		return err
	}

	f.m.Lock()
	f.services[s.ID] = s
	f.m.Unlock()

	return nil
}

func (f *file) delService(id string) {
	f.m.Lock()
	delete(f.services, id)
	f.m.Unlock()
}

func (f *file) listServices() []service {
	f.m.Lock()
	defer f.m.Unlock()

	return f.services.list()
}

func (f *file) getService(id string) (service, bool) {
	f.m.Lock()
	s, ok := f.services[id]
	f.m.Unlock()

	return s, ok
}

func (f *file) checkService(s service) error {
	if s.ID == "" {
		return errors.New("service id cannot be empty")
	}

	if s.Name != f.job {
		return fmt.Errorf("job missmatch: service name %s and file job %s are not equal", s.Name, f.job)
	}

	if s.Namespace != f.namespace {
		return fmt.Errorf("namespace missmatch: service namespace %s and file namespace %s are not equal", s.Namespace, f.namespace)
	}

	return nil
}

func (f *file) data() (data []byte, hash string, err error) {
	svcs := f.listServices()

	if len(svcs) == 0 || f.exportCfg == discovery.Disabled {
		return []byte{}, "", nil
	}

	t := make([]targetGroup, 0, len(f.services))

	for i := range svcs {
		t = append(t, newTargetGroup(svcs[i], f.exportCfg))
	}

	d, err := json.Marshal(t)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal json %s÷ %w", f.relativePath(), err)
	}

	return d, fmt.Sprintf("%x", sha256.Sum256(d)), nil
}

type files struct {
	m               *sync.Mutex
	log             *zap.SugaredLogger
	namespaceGetter namespaceGetter
	files           map[string]*file // files per namespace:jobname
}

func (f files) String() string {
	f.m.Lock()
	defer f.m.Unlock()

	keys := make([]string, 0, len(f.files))
	for k := range f.files {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	buf := bytes.NewBufferString("")

	for _, k := range keys {
		svcs := f.files[k].listServices()
		eps := make([]string, 0, len(svcs))

		for i := range svcs {
			eps = append(eps, svcs[i].Endpoint.String())
		}

		for _, e := range eps {
			fmt.Fprintf(buf, "%s: %s\n", k, e)
		}
	}

	return buf.String()
}

func (f *files) getFiles() map[string]*file {
	f.m.Lock()
	files := f.files
	f.m.Unlock()

	return files
}

func (f *files) reset() {
	f.m.Lock()
	f.files = map[string]*file{}
	f.m.Unlock()
}

func (f *files) addService(s *discovery.Service) error {
	f.m.Lock()
	defer f.m.Unlock()

	svc := service{*s}

	if _, ok := f.files[svc.key()]; !ok {
		ns, err := f.namespaceGetter.Get(svc.Namespace)
		if err != nil {
			return err
		}

		f.files[svc.key()] = &file{
			services:  services{},
			job:       svc.Name,
			namespace: svc.Namespace,
			m:         &sync.Mutex{},
			exportCfg: ns.Export,
		}
	}

	return f.files[svc.key()].addService(svc)
}

// delService deletes service for namespace and id. If service is not found,
// an error is returned.
func (f *files) delService(namespace, id string) error {
	f.m.Lock()
	defer f.m.Unlock()

	for k, file := range f.files {
		if !strings.HasPrefix(k, namespace+":") {
			continue
		}

		if _, ok := file.getService(id); ok {
			file.delService(id)
			return nil
		}
	}

	return errors.New("not found")
}

// getService returns service for namespace and id. If service is not found,
// nil is returned.
func (f *files) getService(namespace, id string) *discovery.Service {
	f.m.Lock()
	defer f.m.Unlock()

	for k, file := range f.files {
		if !strings.HasPrefix(k, namespace+":") {
			continue
		}

		if s, ok := file.getService(id); ok {
			return &s.Service
		}
	}

	return nil
}

// writeFile writes discovery files to destDir. It only writes file if necessary, i.e: if
// there are pending changes.
func (f *files) writeFile(destDir string, file *file) error {
	p := filepath.Join(destDir, file.relativePath())

	data, hash, err := file.data()
	if err != nil {
		return err
	}

	// check for pending changes
	if hash == file.hash {
		return nil
	}

	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}

	f.log.Infow("updating discovery file", "path", p)

	file.hash = hash

	return renameio.WriteFile(p, data, 0600) //nolint: gocritic // we need here the octal value for file permissions.
}

// write writes all file to destDir. It only touches files, that have pending writes.
func (f *files) write(destDir string) error {
	for _, file := range f.getFiles() {
		if file.exportCfg == discovery.Disabled {
			f.log.Debugw("export for namespace is disabled", "namespace", file.namespace)
			continue
		}

		err := f.writeFile(destDir, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *files) delObsoleteFiles(destDir string) error {
	m := map[string]bool{}

	for _, file := range f.getFiles() {
		m[filepath.Join(destDir, file.relativePath())] = true
	}

	return filepath.Walk(destDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !m[path] {
				f.log.Infow("remove obsolete file", "path", path)
				if err := os.Remove(path); err != nil {
					return err
				}
			}

			return nil
		})
}
