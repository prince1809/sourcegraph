package processrestart

import "errors"

// canRestart reports whether the current set of Sourcegraph process can be restarted.
func canRestart() bool {
	return usingGoremanDev || usingGoremanServer
}

// Restart restarts the current set of Sourcegraph processes associated with
// this server.
func Restart() error {
	if !canRestart() {
		return errors.New("reloading site is not supported")
	}
	return nil
}

// willRestart is a channel that is closed when the process will imminently restart.
var WillRestart = make(chan struct{})
