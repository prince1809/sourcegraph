package main

import "github.com/prince1809/sourcegraph/pkg/env"

var (
	FrontendHost = env.Get("LOAD_TEST_FRONTEND_URL", "http://sourcegraph-fronend-internal", "URL to the Sourcegrph frontend host to load test")
	FrontendPort = env.Get("loadTestFrontendPort", "80", "Port that the sourcegraph frontend is listening on")
)
func main()  {
	
}
