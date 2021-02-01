package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/postfinance/discovery"
	"github.com/postfinance/store"
	"k8s.io/apimachinery/pkg/labels"
)

// Common errors
var (
	ErrNotFound = errors.New("not found")
)

// Server represents the server repository.
type Server struct {
	backend store.Backend
	prefix  string
	w       *serverWatcher
}

// NewServer creates a new server repo.
func NewServer(backend store.Backend) *Server {
	return &Server{
		prefix:  serverPrefix,
		backend: backend,
	}
}

// Save saves a saves a new server to backend.
func (s *Server) Save(server discovery.Server) (*discovery.Server, error) {
	server.Modified = time.Now()

	if _, err := store.Put(s.backend, s.key(server.Name), server); err != nil {
		return nil, err
	}

	return &server, nil
}

// Get gets a server by server name.
func (s *Server) Get(serverName string) (*discovery.Server, error) {
	server := discovery.Server{}

	_, err := s.backend.Get(s.key(serverName), store.WithHandler(func(k, v []byte) error {
		err := json.Unmarshal(v, &server)
		if err != nil {
			return err
		}

		return nil
	}))

	if err != nil {
		return nil, err
	}

	return &server, nil
}

// Delete removes given server from backend.
func (s *Server) Delete(name string) error {
	count, err := s.backend.Del(s.key(name))
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("%s: %v", name, ErrNotFound)
	}

	return nil
}

// List lists all servers (if selector is empty string). Otherwise it lists
// the selected servers.
func (s *Server) List(selector string) (discovery.Servers, error) {
	servers := discovery.Servers{}

	_, err := s.backend.Get(s.prefix, store.WithPrefix(), store.WithHandler(func(k, v []byte) error {
		server := discovery.Server{}

		err := json.Unmarshal(v, &server)
		if err != nil {
			return err
		}
		servers = append(servers, server)
		return nil
	}))

	if err != nil {
		return nil, err
	}

	if selector == "" {
		return servers, nil
	}

	sel, err := labels.Parse(selector)
	if err != nil {
		return nil, err
	}

	result := servers.Filter(discovery.ServersBySelector(sel))

	return result, nil
}

func (s *Server) key(serverName string) string {
	return path.Join(s.prefix, serverName)
}

func serverFromKey(key []byte) (*discovery.Server, error) {
	k := string(key)
	fields := strings.Split(k, "/")

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid service path: %s", k)
	}

	name := fields[len(fields)-1]

	return &discovery.Server{
		Name: name,
	}, nil
}
