package main

import (
	"bytes"
	"github.com/gorilla/handlers"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"github.com/prince1809/sourcegraph/pkg/debugserver"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/tracer"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/inconshreveable/log15.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	logRequests, _ = strconv.ParseBool(env.Get("LOG_REQUESTS", "", "log HTTP requests"))
)

const port = "3180"

// requestMu ensures we only do one request at a time to prevent tripping abuse detection.
var requestMu sync.Mutex

var rateLimitingRemainingGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "src",
	Subsystem: "github",
	Name:      "rate_limit_remaining",
	Help:      "Number of calls to Github's API remaining before hitting the rate limit.",
}, []string{"resource"})

func init() {
	rateLimitingRemainingGauge.WithLabelValues("core").Set(5000)
	rateLimitingRemainingGauge.WithLabelValues("search").Set(30)
	prometheus.MustRegister(rateLimitingRemainingGauge)
}

// list obtained from httputil of headers not to forward.
var hopHeaders = map[string]struct{}{
	"Connection":           {},
	"Proxy-Connection":     {}, // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive":           {},
	"Proxy-Authentication": {},
	"Proxy-Authorization":  {},
	"Te":                   {},
	"Trailer":              {},
	"Transfer-Encoding":    {},
	"Upgrade":              {},
}

func main() {
	env.Lock()
	env.HandleHelpFlag()
	tracer.Init()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGHUP)
		<-c
		os.Exit(0)
	}()

	go debugserver.Start()

	var (
		authenticateRequestMu sync.RWMutex
		authenticateRequest   func(query url.Values, header http.Header)
	)
	conf.Watch(func() {
		cfg := conf.Get()
		if clientID, clientSecret := cfg.GithubClientID, cfg.GithubClientSecret; clientID != "" && clientSecret != "" {
			authenticateRequestMu.Lock()
			authenticateRequest = func(query url.Values, header http.Header) {
				query.Set("client_id", clientID)
				query.Set("client_secret", clientSecret)
			}
			authenticateRequestMu.Unlock()
		}
	})

	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q2 := r.URL.Query()
		h2 := make(http.Header)
		for k, v := range r.Header {
			if _, found := hopHeaders[k]; !found {
				h2[k] = v
			}
		}

		// Authenticate for higher rate limits.
		authenticateRequestMu.RLock()
		authRequest := authenticateRequest
		authenticateRequestMu.RUnlock()
		if authRequest != nil {
			authRequest(q2, h2)
		}

		req2 := &http.Request{
			Method: r.Method,
			Body:   r.Body,
			URL: &url.URL{
				Scheme:   "https",
				Host:     "api.github.com",
				Path:     r.URL.Path,
				RawQuery: q2.Encode(),
			},
			Header: h2,
		}

		requestMu.Lock()
		resp, err := http.DefaultClient.Do(req2)
		requestMu.Unlock()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if limit := resp.Header.Get("X-Ratelimit-Remaining"); limit != "" {
			limit, _ := strconv.Atoi(limit)

			resource := "core"
			if strings.HasPrefix(r.URL.Path, "/search/") {
				resource = "search"
			} else if r.URL.Path == "/graphql" {
				resource = "graphql"
			}
			rateLimitingRemainingGauge.WithLabelValues(resource).Set(float64(limit))
		}

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		if resp.StatusCode < 400 || !logRequests {
			io.Copy(w, resp.Body)
			return
		}
		b, err := ioutil.ReadAll(resp.Body)
		log15.Warn("proxy error", "status", resp.StatusCode, "body", string(b), "bodyErr", err)
		io.Copy(w, bytes.NewReader(b))
	})
	if logRequests {
		h = handlers.LoggingHandler(os.Stdout, h)
	}
	h = prometheus.InstrumentHandler("github-proxy", h)
	http.Handle("/", h)

	host := ""
	if env.InsecureDev {
		host = "127.0.0.1"
	}
	addr := net.JoinHostPort(host, port)
	log15.Info("github-proxy: listening", "addr", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
