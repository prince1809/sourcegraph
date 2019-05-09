package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/prince1809/sourcegraph/pkg/env"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
)

var frontendInternal = env.Get("SRC_FRONTEND_INTERNAL", "sourcegraph-frontend-internal", "HTTP address for internal frontend HTTP API.")

type internalClient struct {
	// URL is the root to the internal API frontend server.
	URL string
}

var InternalClient = &internalClient{URL: "http://" + frontendInternal}

// WaitForFrontend retries a noop request to the internal API until it is able to reach
// the endpoint, indicating that the frontend is available.
func (c *internalClient) WaitForFrontend(ctx context.Context) error {
	panic("implement me")
}

// MockInternalClientConfiguration mocks (*internalClient).Configuration.
var MockInternalClientConfiguration func() (conftypes.RawUnified, error)

func (c *internalClient) Configuration(ctx context.Context) (conftypes.RawUnified, error) {
	if MockInternalClientConfiguration != nil {
		return MockInternalClientConfiguration()
	}
	var cfg conftypes.RawUnified
	err := c.postInternal(ctx, "configuration", nil, &cfg)
	return cfg, err
}

// postInternal sends an HTTP post request to the internal route.
func (c *internalClient) postInternal(ctx context.Context, route string, reqBody, respBody interface{}) error {
	return c.post(ctx, "/.internal/"+route, reqBody, respBody)
}

// post sends an HTTP post request to the provided route. If reqBody is
// non-nil it will Marshal it as JSON and set that as the Request body. If
// respBody is non-nil the response body will be JSON unmarshalled to resp.
func (c *internalClient) post(ctx context.Context, route string, reqBody, respBody interface{}) error {
	var data []byte
	if reqBody != nil {
		var err error
		data, err = json.Marshal(reqBody)
		if err != nil {
			return err
		}
	}

	resp, err := ctxhttp.Post(ctx, nil, c.URL+route, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := checkAPIResponse(resp); err != nil {
		return err
	}

	if respBody != nil {
		return json.NewDecoder(resp.Body).Decode(respBody)
	}
	return nil
}

func checkAPIResponse(resp *http.Response) error {
	if 200 > resp.StatusCode || resp.StatusCode > 299 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		b := buf.Bytes()
		errString := string(b)
		if errString != "" {
			return fmt.Errorf("internal API response error code %d: %s (%s)", resp.StatusCode, errString, resp.Request.URL)
		}
		return fmt.Errorf("internal API response error code %d (%s)", resp.StatusCode, resp.Request.URL)
	}
	return nil
}
