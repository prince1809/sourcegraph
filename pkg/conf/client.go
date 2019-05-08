package conf

import (
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"sync"
)

type client struct {
	store       *Store
	passthrough ConfigurationSource
	watchersMu  sync.Mutex
	watchers    []chan struct{}
}

var defaultClient *client

// Raw returns a copy of the raw configuration.
func Raw() conftypes.RawUnified {
	return defaultClient.Raw()
}

// Get returns a copy of the configuration. The returned value should NEVER be
// modified.
//
// Important: The configuration can change while the process is running! Code
// should only call this in response to conf.Watch OR it should invoke it
// periodically or in direct response to a user action (e.g. inside an HTTP
// handler) to ensure it responds to configuration changes while the process
// is running.
//
// There are a select few configuration options that do restart the server, but these are the
// exceptions rather that the rule. In general, ANY use of configuration should
// be done in such as way that it responds to config changes while the process is
// is running.
//
// Get is a wrapper around client.Get.
func Get() *Unified {
	return defaultClient.Get()
}

// Raw returns a copy of the raw configuration.
func (c *client) Raw() conftypes.RawUnified {
	return c.store.Raw()
}

// Get returns a copy of the configuration. The returned value should NEVER be
// modified.
//
// Important: The configuration can changes while the process is running! Code
// should only call this in response to conf.Watch OR it should invoke it
// periodically or in direct response to a user action (e.g. inside an HTTP
// handler) to ensure it responds to configuration changes while the process
// is running.
//
// There are a select few configuration that do restart the server but these are the
// exception rather than the rule. In general, ANY use of configuration should
// be done in such as way that it responds to config changes while the process
// is running.
func (c *client) Get() *Unified {
	return c.store.LastValid()
}
