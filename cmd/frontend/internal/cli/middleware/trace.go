package middleware

import (
	"github.com/prince1809/sourcegraph/pkg/env"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

var httpTrace, _ = strconv.ParseBool(env.Get("HTTP_TRACE", "false", "dump HTTP requests (including body) to stderr"))

// Trace is an HTTP middleware that dumps the HTTP request body (to stderr) if the env var
// `HTTP_TRACE=1`
func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if httpTrace {
			data, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Println("HTTP_TRACE: unable to print request:", err)
			}
			log.Println("==================================================================")
			log.Println(string(data))
			log.Println("=====================================================================")
		}
		next.ServeHTTP(w, r)
	})
}
