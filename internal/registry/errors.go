package registry

import "errors"

// Common errors
var (
	ErrNoServersFound    = errors.New("no servers found")
	ErrNamespaceNotFound = errors.New("namespace not found")
	ErrValidation        = errors.New("validation error")
	ErrContainsServices  = errors.New("server has registered services")
)

// IsServersNotFound returns true on service registration when
// no suitable server is found.
func IsServersNotFound(err error) bool {
	return err == ErrNoServersFound
}

// IsNamespaceNotFound returns true on service registration when
// the specified namespace does not exist.
func IsNamespaceNotFound(err error) bool {
	return err == ErrNoServersFound
}

// IsValidationError returns true if a validation error occurred.
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}
