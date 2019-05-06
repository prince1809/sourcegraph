package conf

import (
	"context"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"sync"
)

// ConfigurationSource provides direct access to read and write to the
// "raw" configuration.
type ConfigurationSource interface {
	// Write updates the configuration. The Deployment field is ignored.
	Write(ctx context.Context, data conftypes.RawUnified) error
	Read(ctx context.Context) (conftypes.RawUnified, error)
}

type Server struct {
	Source ConfigurationSource

	// fileWrite signals when our app writes to the configuration file. The
	// secondary channel is closed when server.Raw() would return the new
	// configuration that has been written to disk.
	fileWrite chan chan struct{}

	once sync.Once
}

// NewServer returns a new Server instance that manages the site config file
// that is stored at configSource.
//
// The server must be started with Start() before it can handle requests.
func NewServer(source ConfigurationSource) *Server {
	fileWrite := make(chan chan struct{}, 1)
	return &Server{
		Source:    source,
		fileWrite: fileWrite,
	}
}

// Write writes the JSON config file to the config file's path. If the JSON configuration is
// invalid, an error is returned.
func (s *Server) Write(ctx context.Context, input conftypes.RawUnified) error {
	// Parse the configuration so that we can diff it (this also validates it
	// is proper JSON).
	panic("implement me")
	return nil
}
