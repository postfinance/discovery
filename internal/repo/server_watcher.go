package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/postfinance/store"
	"github.com/zbindenren/discovery"
)

// ServerEvent contains the server and the event (change or delete).
type ServerEvent struct {
	discovery.Server
	Event Event
}

// Chan returns read-only channel of server events.
//nolint: dupl // difficult to dedupl
func (s *Server) Chan(ctx context.Context, errorHandler func(error)) <-chan *ServerEvent {
	if s.w != nil {
		return s.w.c
	}

	c := make(chan *ServerEvent, 1)
	s.w = &serverWatcher{
		c: c,
	}

	watchReady := make(chan struct{})
	notifyCreated := func() {
		close(watchReady)
	}

	go func() {
		if err := s.backend.Watch(s.prefix, s.w,
			store.WithContext(ctx),
			store.WithNotifyCreated(notifyCreated),
			store.WithPrefix(),
		); err != nil {
			errorHandler(err)
		}
	}()

	<-watchReady

	return c
}

type serverWatcher struct {
	c chan *ServerEvent
}

// OnPut implements Watcher interface.
func (w serverWatcher) OnPut(k, v []byte) error {
	e, err := w.onChange(k, v)
	if err != nil {
		return err
	}

	w.c <- e

	return nil
}

// OnDelete implements Watcher interface.
func (w serverWatcher) OnDelete(k, v []byte) error {
	e, err := w.onDelete(k)
	if err != nil {
		return err
	}

	w.c <- e

	return nil
}

// BeforeWatch implements Watcher interface.
func (w serverWatcher) BeforeWatch() error {
	return nil
}

// BeforeLoop implements Watcher interface.
func (w serverWatcher) BeforeLoop() error {
	return nil
}

// OnDone implements Watcher interface.
func (w serverWatcher) OnDone() error {
	close(w.c)

	return nil
}

func (w serverWatcher) onChange(k, v []byte) (*ServerEvent, error) {
	s := discovery.Server{}

	if err := json.Unmarshal(v, &s); err != nil {
		return nil, fmt.Errorf("failed to unmashal server %s: %w", string(k), err)
	}

	return &ServerEvent{
		Event:  Change,
		Server: s,
	}, nil
}

func (w serverWatcher) onDelete(k []byte) (*ServerEvent, error) {
	s, err := serverFromKey(k)
	if err != nil {
		return nil, err
	}

	return &ServerEvent{
		Event:  Delete,
		Server: *s,
	}, nil
}
