package conf

import (
	"context"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/pkg/api"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"sync"
	"time"
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

// Watch calls the given function whenever the configuration has changed. The new configuration is
// accessed by calling conf.Get.
//
// Before Watch returns, it will invoke f to use the current configuration.
//
// Watch is a wrapper around client.Watch.
//
// IMPORTANT: Watch will block on config initialization. It therefore should *never* be called
// synchronously in `init` functions.
func Watch(f func()) {
	defaultClient.Watch(f)
}

// Watch calls the given function in a separate goroutine whenever the
// configuration has changed. The new configuration can be received by callling
// conf.Get.
// Before Watch returns, it will invoke f to use the current configuration.
func (c *client) Watch(f func()) {
	// Add the watcher channel now, rather than after invoking f below, in case
	// an update were to happen while we were invoking f.
	notify := make(chan struct{}, 1)
	c.watchersMu.Lock()
	c.watchers = append(c.watchers, notify)
	c.watchersMu.Unlock()

	// Call the function now, to use the current configuration.
	c.store.WaitUntilInitialized()
	f()

	go func() {
		// Invoke f when the configuration has changed.
		for {
			<-notify
			f()
		}
	}()
}

type continuousUpdateOptions struct {
	// delayBeforeUnreachableLog is how long to wait before logging an error upon initial startup
	// due to the frontend being unreachable. It is used to avoid log spam when other services (that
	// contact the frontend for configuration) start up before the frontend.
	delayBeforeUnreachableLog time.Duration

	log   func(format string, v ...interface{}) // log.Printf equivalent
	sleep func()                                // sleep between updates
}

// continuouslyUpdate runs (*client).fetchAndUpdate in an infinite loop, with error logging and
// random sleep intervals.
//
// The optOnlySetByTests parameter is ONLY customized by tests. ALl callers in main code should pass
// nil (so that the same defaults are used).
func (c *client) continuouslyUpdate(optOnlySetByTests *continuousUpdateOptions) {

}

func (c *client) fetchAndUpdate() error {
	ctx := context.Background()
	var (
		newConfig conftypes.RawUnified
		err       error
	)
	if c.passthrough != nil {
		newConfig, err = c.passthrough.Read(ctx)
	} else {
		newConfig, err = api.InternalClient.Configuration(ctx)
	}
	if err != nil {
		return errors.Wrap(err, "unable to fetch new configuration")
	}

	configChange, err := c.store.MaybeUpdate(newConfig)
	if err != nil {
		return errors.Wrap(err, "unable to update new configuration")
	}

	if configChange.Changed {
		//c.NotifyWatchers()
	}
	return nil

}
