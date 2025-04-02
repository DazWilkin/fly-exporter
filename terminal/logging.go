// Replaced flyctl/terminal with this implementation
// flyctl/terminal takes a dependency on flyctl/internal
// And this caused Modules issues requiring the use of a replace in go.mod:
// github.com/loadsmart/calver-go => github.com/ndarilek/calver-go v0.0.0-20230710153822-893bbd83a936

package terminal

import (
	"encoding/json"
	"fmt"
	"io"

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
func (t Terminal) Debug(v ...any) {
	logContent := func(b []byte) {
		if len(b) == 0 {
			return
		}

		if json.Valid(b) {
			var j map[string]any
			if err := json.Unmarshal(b, &j); err != nil {
				t.log.Info("debug",
					"err", err,
				)
				return
			}

			t.log.Info("debug",
				"json", j,
			)
			return
		}

		t.log.Info("debug",
			"string", string(b),
		)
	}

	for _, item := range v {
		switch typedItem := item.(type) {
		case string:
			b := []byte(typedItem)
			logContent(b)
		case io.Reader:
			b, err := io.ReadAll(typedItem)
			if err != nil {
				t.log.Info("debug", "err", err)
				continue
			}
			logContent(b)
		default:
			t.log.Info("debug",
				"type", fmt.Sprintf("%T", typedItem),
				"content", typedItem,
			)
		}
	}
}

// Debugf is a method that wraps go-logr/logger.Info
func (t Terminal) Debugf(format string, v ...any) {
	t.log.Info(format, v...)
}
