package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/postfinance/store"
	"github.com/zbindenren/discovery"
)

// Event represents a repo event. The possible valid values are Change and Delete.
type Event int

// All known events. An Unknown event should never occur.
const (
	Unknown Event = iota
	Change
	Delete
)

// ServiceEvent contains the service and the event (change or delete).
type ServiceEvent struct {
	discovery.Service
	Event Event
}

// Chan returns read-only channel of service events.
//nolint: dupl // difficult to dedupl
func (s *Service) Chan(ctx context.Context, errorHandler func(error)) <-chan *ServiceEvent {
	if s.w != nil {
		return s.w.c
	}

	c := make(chan *ServiceEvent, 1)
	s.w = &serviceWatcher{
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

type serviceWatcher struct {
	c chan *ServiceEvent
}

// OnPut implements Watcher interface.
func (w serviceWatcher) OnPut(k, v []byte) error {
	e, err := w.onChange(k, v)
	if err != nil {
		return err
	}

	w.c <- e

	return nil
}

// OnDelete implements Watcher interface.
func (w serviceWatcher) OnDelete(k, v []byte) error {
	e, err := w.onDelete(k)
	if err != nil {
		return err
	}

	w.c <- e

	return nil
}

// BeforeWatch implements Watcher interface.
func (w serviceWatcher) BeforeWatch() error {
	return nil
}

// BeforeLoop implements Watcher interface.
func (w serviceWatcher) BeforeLoop() error {
	return nil
}

// OnDone implements Watcher interface.
func (w serviceWatcher) OnDone() error {
	close(w.c)

	return nil
}

func (w serviceWatcher) onChange(k, v []byte) (*ServiceEvent, error) {
	s := discovery.Service{}

	if err := json.Unmarshal(v, &s); err != nil {
		return nil, fmt.Errorf("failed to unmashal service %s: %w", string(k), err)
	}

	return &ServiceEvent{
		Event:   Change,
		Service: s,
	}, nil
}

func (w serviceWatcher) onDelete(k []byte) (*ServiceEvent, error) {
	s, err := serviceFromKey(k)
	if err != nil {
		return nil, err
	}

	return &ServiceEvent{
		Event:   Delete,
		Service: *s,
	}, nil
}
