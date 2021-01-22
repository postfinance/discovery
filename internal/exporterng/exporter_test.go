// +build bla

package exporter

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/postfinance/flash"
	"github.com/postfinance/store/hash"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zbindenren/discovery"
	"github.com/zbindenren/discovery/internal/registry"
	"github.com/zbindenren/discovery/internal/repo"
)

const (
	server = "server1"
)

func TestExporter(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "discovery")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	b, err := hash.New()
	require.NoError(t, err)

	l := flash.New().Get()

	e := New(b, l, dir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reg, err := registry.New(b, prometheus.NewRegistry(), l, 1)
	require.NoError(t, err)

	expected := expected{}
	_, err = reg.RegisterService(*discovery.MustNewService("test1", "http://0.example.com"))
	require.NoError(t, err)
	expected.add("test1", "0.example.com")

	err = ioutil.WriteFile(path.Join(dir, "do-be-removed.json"), []byte{}, 0600)
	require.NoError(t, err)

	go func() {
		err := e.Start(ctx, 24*time.Hour)
		assert.NoError(t, err)
	}()

	_, err = reg.RegisterServer(server, discovery.Labels{})
	require.NoError(t, err)

	t.Run("registering services", func(t *testing.T) {
		_, err = reg.RegisterService(*discovery.MustNewService("test1", "http://1.example.com"))
		require.NoError(t, err)
		expected.add("test1", "1.example.com")

		_, err = reg.RegisterService(*discovery.MustNewService("test1", "http://2.example.com"))
		require.NoError(t, err)

		_, err = reg.RegisterService(*discovery.MustNewService("test2", "http://3.example.com"))
		require.NoError(t, err)

		_, err = reg.RegisterService(*discovery.MustNewService("test3", "http://4.example.com"))
		require.NoError(t, err)
	})

	t.Run("changing name", func(t *testing.T) {
		_, err = reg.RegisterService(*discovery.MustNewService("test2", "http://2.example.com"))
		require.NoError(t, err)
		expected.add("test2", "2.example.com")
	})

	t.Run("unregistering services", func(t *testing.T) {
		err := reg.UnRegisterService("http://4.example.com", "")
		require.NoError(t, err)

		err = reg.UnRegisterService("http://3.example.com", "")
		require.NoError(t, err)
	})

	t.Run("registering last service", func(t *testing.T) {
		_, err = reg.RegisterService(*discovery.MustNewService("last", "http://example.com"))
		require.NoError(t, err)
		expected.add("last", "example.com")
	})

	time.Sleep(2 * time.Second)
	expected.checkAll(t, dir)
}

func TestExporter2(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "discovery")
	require.NoError(t, err)

	ch := make(chan *repo.ServiceEvent)
	mock := newMock(ch)

	l := log.New()
	l.SetDebug(true)

	e := Exporter{
		log:          l.Get(),
		repo:         mock,
		directory:    dir,
		destinations: files{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err = e.Start(ctx, 1*time.Hour)
		require.NoError(t, err)
	}()

	ch <- mock.newEvent(repo.Change, "test1", "www1")
	ch <- mock.newEvent(repo.Change, "test1", "www2")

	time.Sleep(1 * time.Minute)
}

type serviceRepoMock struct {
	servers         discovery.Servers
	ch              chan *repo.ServiceEvent
	initialServices discovery.Services
}

func newMock(ch chan *repo.ServiceEvent) *serviceRepoMock {
	mock := &serviceRepoMock{
		ch: ch,
		servers: discovery.Servers{
			discovery.Server{
				Name: "server1",
			},
		},
	}

	se := mock.newEvent(repo.Change, "init", "init")
	mock.initialServices = discovery.Services{
		se.Service,
	}

	return mock
}

func (s *serviceRepoMock) newEvent(event repo.Event, job, host string) *repo.ServiceEvent {
	svc := discovery.MustNewService(job, fmt.Sprintf("http://%s.example.com", host))
	svc.Servers = s.servers.Names()

	return &repo.ServiceEvent{
		Event:   event,
		Service: *svc,
	}
}

func (s *serviceRepoMock) List(namespace, selector string) (discovery.Services, error) {
	return s.initialServices, nil
}

func (s *serviceRepoMock) Chan(context.Context, func(error)) <-chan *repo.ServiceEvent {
	return s.ch
}

type expected map[string][]string

func (e expected) add(job, endpoint string) {
	e[job] = append(e[job], endpoint)
}

func (e expected) ok(dir, job string) error {
	absPath := createAbsPath(dir, server, "default", job)

	d, err := ioutil.ReadFile(absPath) // nolint: gosec
	if err != nil {
		return err
	}

	for _, ep := range e[job] {
		if !bytes.Contains(d, []byte(ep)) {
			return fmt.Errorf("endoint %s not found in %s", ep, absPath)
		}
	}

	return nil
}

func (e expected) checkAll(t *testing.T, dir string) {
	for job := range e {
		err := e.ok(dir, job)
		require.NoError(t, err)
	}
}
