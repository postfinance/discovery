package repo

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/postfinance/discovery"
	"github.com/postfinance/store"
)

// Namespace represents the server repository.
type Namespace struct {
	backend store.Backend
	prefix  string
	w       *namespaceWatcher
}

// NewNamespace creates a new namespace repo.
func NewNamespace(backend store.Backend) *Namespace {
	return &Namespace{
		prefix:  namespacePrefix,
		backend: backend,
	}
}

// Save saves a saves a new server to backend.
func (n *Namespace) Save(namespace discovery.Namespace) (*discovery.Namespace, error) {
	namespace.Modified = time.Now()

	if _, err := store.Put(n.backend, n.key(namespace.Name), namespace); err != nil {
		return nil, err
	}

	return &namespace, nil
}

// Get gets a server by server name.
func (n *Namespace) Get(namespaceName string) (*discovery.Namespace, error) {
	namespace := discovery.Namespace{}

	_, err := n.backend.Get(n.key(namespaceName), store.WithHandler(func(k, v []byte) error {
		err := json.Unmarshal(v, &namespace)
		if err != nil {
			return err
		}

		return nil
	}))

	if err != nil {
		return nil, err
	}

	return &namespace, nil
}

// Delete removes given server from backend.
func (n *Namespace) Delete(name string) error {
	count, err := n.backend.Del(n.key(name))
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return nil
}

// List lists all servers (if selector is empty string). Otherwise it lists
// the selected servers.
func (n *Namespace) List() (discovery.Namespaces, error) {
	namespaces := discovery.Namespaces{}

	_, err := n.backend.Get(n.prefix, store.WithPrefix(), store.WithHandler(func(k, v []byte) error {
		namespace := discovery.Namespace{}

		err := json.Unmarshal(v, &namespace)
		if err != nil {
			return err
		}
		namespaces = append(namespaces, namespace)
		return nil
	}))

	if err != nil {
		return nil, err
	}

	return namespaces, nil
}

func (n *Namespace) key(namespaceName string) string {
	return path.Join(n.prefix, namespaceName)
}

func namespaceFromKey(key []byte) (*discovery.Namespace, error) {
	k := string(key)
	fields := strings.Split(k, "/")

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid service path: %s", k)
	}

	name := fields[len(fields)-1]

	return &discovery.Namespace{
		Name: name,
	}, nil
}
