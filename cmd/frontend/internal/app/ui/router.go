package ui

import (
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	uirouter "github.com/prince1809/sourcegraph/cmd/frontend/internal/app/ui/router"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/randstring"
	"github.com/prince1809/sourcegraph/pkg/trace"
	"gopkg.in/inconshreveable/log15.v2"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
)

const (
	routeHome           = "home"
	routeStart          = "start"
	routeSearch         = "search"
	routeSearchBadge    = "search-badge"
	routeSearchSearches = "search-searches"

	aboutRedirectScheme = "https"
	aboutRedirectHost   = "about.sourcegraph.com"

	routeWelcome = "welcome"
)

// Router returns the router that serves pages for our web app.
func Router() *mux.Router {
	return uirouter.Router
}

var (
	mockServeRepo func(w http.ResponseWriter, r *http.Request)
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	// Top-level routes.
	r.Path("/").Methods("GET").Name(routeHome)
	r.Path("/start").Methods("GET").Name(routeStart)
	r.PathPrefix("/welcome").Methods("GET").Name(routeWelcome)
	r.Path("/search").Methods("GET").Name(routeSearch)

	return r
}

func init() {
	initRouter()
}

func initRouter() {
	// basic pages with static titles
	router := newRouter()
	uirouter.Router = router // make accessible to other packages
	router.Get(routeHome).Handler(handler(serveHome))
	router.Get(routeStart).Handler(staticRedirectHandler("/welcome", http.StatusMovedPermanently))
	router.Get(routeWelcome).Handler(handler(serveWelcome))

	// search
	router.Get(routeSearch).Handler(handler(serveBasicPage(func(c *Common, r *http.Request) string {
		shortQuery := limitString(r.URL.Query().Get("q"), 25, true)
		if shortQuery == "" {
			return "Sourcegraph"
		}

		return fmt.Sprintf("%s - Sourcegraph", shortQuery)
	})))
}

// staticRedirectHandler returns an HTTP handler that redirects all requests to
// the specified url or relative path with the specified status code.
//
// The scheme, host and path in the specified url overrides ones in the incoming
// request. For example:
//
// staticRedirectHandler("http://google.com") serving "https://sourcegraph.com/foobar?q=foo" -> "https://google.com/foobar?q=foo"
// staticRedirectHandler("/foo") serving "https://"
func staticRedirectHandler(u string, code int) http.Handler {
	target, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if target.Scheme != "" {
			r.URL.Scheme = target.Scheme
		}
		if target.Host != "" {
			r.URL.Host = target.Host
		}
		if target.Path != "" {
			r.URL.Path = target.Path
		}
		http.Redirect(w, r, r.URL.String(), code)
	})
}

func limitString(s string, n int, ellipsis bool) string {
	if len(s) < n {
		return s
	}
	if ellipsis {
		return s[:n-1] + "…"
	}
	return s[:n-1]
}

// handler wraps an HTTP handler that returns potential errors. If any error is
// returned, serveError is called.
//
// Clients that wish to return their own HTTP status code should use this from
// their handler:
//
// serveError(w,r, err, http.MyStatusCode)
// return nil
//
func handler(f func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				serveError(w, r, recoverError{recover: rec, stack: debug.Stack()}, http.StatusInternalServerError)
			}
		}()
		if err := f(w, r); err != nil {
			serveError(w, r, err, http.StatusInternalServerError)
		}
	})
	return trace.TraceRoute(gziphandler.GzipHandler(h))
}

type recoverError struct {
	recover interface{}
	stack   []byte
}

func (r recoverError) Error() string {
	return fmt.Sprintf("ui: recovered from panic: %v", r.recover)
}

// serveError serves the error template with the specified error message. It is
// assumed that the error message could accidentally contain sensitive data,
// and as such is only presented to the user in debug mode.
func serveError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	serveErrorNoDebug(w, r, err, statusCode, false, false)
}

type pageError struct {
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
	Error      string `json:"error"`
	ErrorID    string `json:"errorID"`
}

// serveErrorNoDebug should not be called by anyone except serveErrorTest.
func serveErrorNoDebug(w http.ResponseWriter, r *http.Request, err error, statusCode int, nodebug, forceServeError bool) {
	w.WriteHeader(statusCode)
	errorID := randstring.NewLen(16)

	// Determine span URL and log the error.
	var spanURL string
	if span := opentracing.SpanFromContext(r.Context()); span != nil {

	}
	log15.Error("ui HTTP handler error response", "method", r.Method, "request_uri", r.URL.RequestURI(), "statuc_code", statusCode, "error", err, "error_id", errorID, "trace", spanURL)

	// In the case of recovering from a panic, we nicely include the stack
	// trace in the error that is shown on the page. Additionally, we log it
	// separately (since log15 prints the escaped sequence).
	if r, ok := err.(recoverError); ok {
		err = fmt.Errorf("ui: recovered from panic %v\n\n%s", r.recover, r.stack)
		log.Println(err)
	}

	var errorIfDebug string
	if forceServeError || (env.InsecureDev && !nodebug) {
		errorIfDebug = err.Error()
	}

	pageErrorContext := &pageError{
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
		Error:      errorIfDebug,
		ErrorID:    errorID,
	}

	// First try to render the error fancily: this relies on *Common
	// functionality that might always work (for example, if some services are
	// down rather than something that is primarily a user error).
	delete(mux.Vars(r), "Repo")
	var commonServeErr error
	title := fmt.Sprintf("%v %s - Sourcegraph", statusCode, http.StatusText(statusCode))
	common, commonErr := newCommon(w, r, title, func(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
		// Stub our serveErrorContext to newCommon so that it is not reentrant.
		commonServeErr = err
	})
	common.Error = pageErrorContext
	if commonErr == nil && commonServeErr == nil {

	}

	// Fallback to ugly/ reliable error template
	stdErr := renderTemplate(w, "error.html", pageErrorContext)
	if stdErr != nil {
		log15.Error("ui: error while serving final error template", "error", stdErr)
	}
}
