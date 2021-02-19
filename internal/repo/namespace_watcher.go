package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/postfinance/discovery"
	"github.com/postfinance/store"
)

// NamespaceEvent contains the server and the event (change or delete).
type NamespaceEvent struct {
	discovery.Namespace
	Event Event
}

// Chan returns read-only channel of server events.
func (n *Namespace) Chan(ctx context.Context, errorHandler func(error)) <-chan *NamespaceEvent {
	if n.w != nil {
		return n.w.c
	}

	c := make(chan *NamespaceEvent, 1)
	n.w = &namespaceWatcher{
		c: c,
	}

	watchReady := make(chan struct{})
	notifyCreated := func() {
		close(watchReady)
	}

	go func() {
		if err := n.backend.Watch(n.prefix, n.w,
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

type namespaceWatcher struct {
	c chan *NamespaceEvent
}

// OnPut implements Watcher interface.
func (w namespaceWatcher) OnPut(k, v []byte) error {
	e, err := w.onChange(k, v)
	if err != nil {
		return err
	}

	w.c <- e

	return nil
}

// OnDelete implements Watcher interface.
func (w namespaceWatcher) OnDelete(k, v []byte) error {
	e, err := w.onDelete(k)
	if err != nil {
		return err
	}

	w.c <- e

	return nil
}

// BeforeWatch implements Watcher interface.
func (w namespaceWatcher) BeforeWatch() error {
	return nil
}

// BeforeLoop implements Watcher interface.
func (w namespaceWatcher) BeforeLoop() error {
	return nil
}

// OnDone implements Watcher interface.
func (w namespaceWatcher) OnDone() error {
	close(w.c)

	return nil
}

func (w namespaceWatcher) onChange(k, v []byte) (*NamespaceEvent, error) {
	n := discovery.Namespace{}

	if err := json.Unmarshal(v, &n); err != nil {
		return nil, fmt.Errorf("failed to unmashal namespace %s: %w", string(k), err)
	}

	return &NamespaceEvent{
		Event:     Change,
		Namespace: n,
	}, nil
}

func (w namespaceWatcher) onDelete(k []byte) (*NamespaceEvent, error) {
	n, err := namespaceFromKey(k)
	if err != nil {
		return nil, err
	}

	return &NamespaceEvent{
		Event:     Delete,
		Namespace: *n,
	}, nil
}
