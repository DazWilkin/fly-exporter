// Replaced flyctl/terminal with this implementation
// flyctl/terminal takes a dependency on flyctl/internal
// And this caused Modules issues requiring the use of a replace in go.mod:
// github.com/loadsmart/calver-go => github.com/ndarilek/calver-go v0.0.0-20230710153822-893bbd83a936

package terminal

import (
	"github.com/go-logr/logr"

	"github.com/superfly/flyctl/api"
)

// Only requires 2 methods to be implemented: Debug, Debugf
// https://pkg.go.dev/github.com/superfly/flyctl/api#Logger
// Even though flyctl/terminal implementation supports many
// https://pkg.go.dev/github.com/superfly/flyctl@v0.1.65/terminal
var _ api.Logger = (*Terminal)(nil)

// Terminal is a struct that replicates flyctl terminal
type Terminal struct {
	log logr.Logger
}

// New wraps a go-logr/logger in a Terminal
func New(log logr.Logger) Terminal {
	return Terminal{
		log: log,
	}
}

// Debug is a method that wraps go-logr/logger.Info
func (t Terminal) Debug(v ...interface{}) {
	t.log.Info("debug", v...)
}

// Debugf is a method that wraps go-logr/logger.Info
func (t Terminal) Debugf(format string, v ...interface{}) {
	t.log.Info(format, v...)
}
