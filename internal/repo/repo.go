// Package repo is responsilbe to store services and servers
// in a (etcd) repository.
package repo

// prefixes relative to the global prefix
const (
	namespacePrefix = "namespace/v1"
	serverPrefix    = "server/v1"
	servicePrefix   = "service/v1"
)
