package exporter

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/postfinance/flash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/repo"
)

func TestEnableDisableWatch(t *testing.T) {
	e := Exporter{}
	e.disableWatch()
	assert.True(t, e.isWatchDisabled())
	e.enableWatch()
	assert.False(t, e.isWatchDisabled())
}

func TestExporter(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "discovery")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	ch := make(chan *repo.ServiceEvent)
	serviceGetter := newServiceMock(ch)
	serverGetter := newServerMock()
	l := flash.New()
	l.SetDebug(true)
	e := Exporter{
		log:         l.Get(),
		serverRepo:  serverGetter,
		serviceRepo: serviceGetter,
		destDir:     dir,
		destinations: files{
			m:               &sync.Mutex{},
			files:           map[string]*file{},
			log:             l.Get(),
			namespaceGetter: newNamespaceMock(),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = e.Start(ctx, "not-exist", 24*time.Hour)
	require.Error(t, err)
	go func() {
		err = e.Start(ctx, "server1", 24*time.Hour)
		require.NoError(t, err)
	}()

	assertFileContains(t, filepath.Join(dir, "server1/standard/default/initial.json"), "initial1.pnet.ch")
	assertFileContains(t, filepath.Join(dir, "server1/standard/default/initial.json"), "initial2.pnet.ch")
	assertFileNotContains(t, filepath.Join(dir, "server1/standard/default/initial.json"), "initial3.pnet.ch")

	serviceGetter.addEvent(&repo.ServiceEvent{
		Event:   repo.Change,
		Service: newService("i1", "changedjob", "https://initial1.pnet.ch"),
	})

	assertFileContains(t, filepath.Join(dir, "server1/standard/default/initial.json"), "initial2.pnet.ch")
	assertFileNotContains(t, filepath.Join(dir, "server1/standard/default/initial.json"), "initial1.pnet.ch")
	assertFileContains(t, filepath.Join(dir, "server1/standard/default/changedjob.json"), "initial1.pnet.ch")
	assertFileNotContains(t, filepath.Join(dir, "server1/standard/default/changedjob.json"), "initial2.pnet.ch")

	serviceGetter.addEvent(&repo.ServiceEvent{
		Event:   repo.Change,
		Service: newService("d1", "initial", "https://initial33.pnet.ch", "other-server1"),
	})
	assertFileNotContains(t, filepath.Join(dir, "server1/standard/default/initial.json"), "initial33.pnet.ch")

	serviceGetter.addEvent(&repo.ServiceEvent{
		Event:   repo.Delete,
		Service: newService("i1", "changedjob", "https://initial1.pnet.ch"),
	})

	assertFileNotContains(t, filepath.Join(dir, "server1/standard/default/changedjob.json"), "initial") // empty

	blackboxSvc := newService("b1", "blackbox", "https://blackbox1.pnet.ch")
	blackboxSvc.Namespace = "appl-blackbox"
	serviceGetter.addEvent(&repo.ServiceEvent{
		Event:   repo.Change,
		Service: blackboxSvc,
	})

	assertFileContains(t, filepath.Join(dir, "server1/blackbox/appl-blackbox/blackbox.json"), `[{"targets":["https://blackbox1.pnet.ch"]}]`)
}

type serviceRepoMock struct {
	ch              chan *repo.ServiceEvent
	initialServices map[string]discovery.Service
}

func newServiceMock(ch chan *repo.ServiceEvent) *serviceRepoMock {
	mock := &serviceRepoMock{
		ch: ch,
		initialServices: map[string]discovery.Service{
			"i1": newService("i1", "initial", "https://initial1.pnet.ch"),
			"i2": newService("i2", "initial", "https://initial2.pnet.ch"),
			"o1": newService("o1", "otherServer", "https://initial3.pnet.ch", "other-server1"),
		},
	}

	return mock
}

func (s *serviceRepoMock) addEvent(se *repo.ServiceEvent) {
	_, ok := s.initialServices[se.Service.ID]

	if ok && se.Event == repo.Delete {
		delete(s.initialServices, se.Service.ID)
		s.ch <- se

		return
	}

	s.initialServices[se.Service.ID] = se.Service

	s.ch <- se
}

func (s *serviceRepoMock) List(namespace, selector string) (discovery.Services, error) {
	services := discovery.Services{}
	for _, svc := range s.initialServices {
		services = append(services, svc)
	}

	return services, nil
}

func (s *serviceRepoMock) Chan(context.Context, func(error)) <-chan *repo.ServiceEvent {
	return s.ch
}

type serverRepoMock struct {
	servers map[string]*discovery.Server
	ch      chan *repo.ServerEvent
}

func (s *serverRepoMock) Get(serverName string) (*discovery.Server, error) {
	server, ok := s.servers[serverName]
	if !ok {
		return nil, errors.New("not found")
	}

	return server, nil
}

func (s *serverRepoMock) Chan(context.Context, func(error)) <-chan *repo.ServerEvent {
	return s.ch
}

func newServerMock() *serverRepoMock {
	mock := &serverRepoMock{
		servers: map[string]*discovery.Server{
			"server1":       discovery.NewServer("server1", discovery.Labels{}),
			"other-server1": discovery.NewServer("other-server1", discovery.Labels{}),
		},
	}

	return mock
}

type namespaceRepoMock struct {
	namespaces map[string]*discovery.Namespace
}

func (n *namespaceRepoMock) Get(namespace string) (*discovery.Namespace, error) {
	ns, ok := n.namespaces[namespace]
	if !ok {
		return nil, errors.New("not found")
	}

	return ns, nil
}

func newNamespaceMock() *namespaceRepoMock {
	mock := &namespaceRepoMock{
		namespaces: map[string]*discovery.Namespace{
			"default": &discovery.Namespace{
				Name:     "default",
				Export:   discovery.Standard,
				Modified: time.Now(),
			},
			"appl-blackbox": &discovery.Namespace{
				Name:     "appl-blackbox",
				Export:   discovery.Blackbox,
				Modified: time.Now(),
			},
		},
	}

	return mock
}

func newService(id, name, endpoint string, servers ...string) discovery.Service {
	s := *discovery.MustNewService(name, endpoint)
	s.ID = id

	if len(servers) == 0 {
		s.Servers = []string{"server1"}
	}

	return s
}

func assertFileContains(t *testing.T, path, substring string) {
	require.Eventually(t, func() bool {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}

		d, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		return strings.Contains(string(d), substring)
	}, time.Second, 10*time.Millisecond)
}

func assertFileNotContains(t *testing.T, path, substring string) {
	require.Eventually(t, func() bool {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}

		d, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		return !strings.Contains(string(d), substring)
	}, time.Second, 10*time.Millisecond)
}
