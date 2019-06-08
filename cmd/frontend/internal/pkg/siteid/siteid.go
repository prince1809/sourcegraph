// Package siteid provides access to the site ID, a stable identifier for the current
// Sourcegraph site.
//
// All servers that are part of the same logical Sourcegraph site have the same site ID
// (although it is possible for an admin to misconfigure the servers so that this in not
// true).
//
// The "site ID" was formally known as the "app ID".
package siteid

import (
	"context"
	"log"
	"os"
	"time"
)

var (
	inited bool
	siteID string

	fatalln = log.Fatalln // overridden in tests
)

// Init reads (or generates) the site ID. This func must be called exactly once before
// Get can be called.
func Init() {
	if inited {
		panic("siteid: already initialised")
	}

	if v := os.Getenv("TRACKING_APP_ID"); v != "" {
		siteID = v
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		globalState, err := globalstatedb.Get(ctx)
		if err != nil {
			fatalln("Error initializing global state", err)
		}
		siteID = globalState.SiteID
	}
	inited = true
}
