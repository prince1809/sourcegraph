package shared

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/cmd/management-console/assets"
	"github.com/prince1809/sourcegraph/cmd/management-console/shared/internal/tlscertgen"
	"github.com/prince1809/sourcegraph/pkg/db/confdb"
	"github.com/prince1809/sourcegraph/pkg/db/dbconn"
	"github.com/prince1809/sourcegraph/pkg/debugserver"
	"github.com/prince1809/sourcegraph/pkg/env"
	"gopkg.in/inconshreveable/log15.v2"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const port = "2633"

var (
	tlsCert       = env.Get("TLS_CERT", "/etc/sourcegraph/management/cert.pem", "TLS certificate (automatically generated if file does not exist)")
	tlsKey        = env.Get("TLS_KEY", "/etc/sourcegraph/management/key.pem", "TLS key (automatically generated if file does not exist)")
	customTLS     = env.Get("CUSTOM_TLS", "false", "When true, disables TLS cert/key generation to prevent accidents.")
	unsafeNoHTTPS = env.Get("UNSAFE_NO_HTTPS", "false", "(unsafe) When true, disables HTTPS entirely. Anyone who can MITM your traffic to the management console can steal the admin password and act on your behalf!")
)

func configureTLS() error {
	customTLS, _ := strconv.ParseBool(customTLS)

	generate := false
	_, err := os.Stat(tlsCert)
	if os.IsNotExist(err) {
		if customTLS {
			return err
		}
		generate = true
	} else if err != nil {
		return err
	}

	_, err = os.Stat(tlsKey)
	if os.IsNotExist(err) {
		if customTLS {
			return err
		}
		generate = true
	} else if err != nil {
		return err
	}

	if customTLS || !generate {
		return nil
	}
	log.Println("Generating and using self-signed TLS cert/key")

	if err := os.MkdirAll(filepath.Dir(tlsCert), 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(tlsKey), 0700); err != nil {
		return err
	}

	// Generate a TLS cert.
	certOut, err := os.Create(tlsCert)
	if err != nil {
		return errors.Wrap(err, "failed to open cert.pem for writing")
	}
	defer certOut.Close()

	keyOut, err := os.OpenFile(tlsKey, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to open key.pem for writing")
	}
	defer keyOut.Close()

	return tlscertgen.Generate(tlscertgen.Options{
		Cert:         certOut,
		Key:          keyOut,
		Hosts:        []string{"management-console.sourcegraph.com"},
		ValidFor:     100 * 365 * 24 * time.Hour,
		ECDSACurve:   "P256",
		Organization: "Sourcegraph.com",
	})

}

func Main() {
	env.Lock()
	env.HandleHelpFlag()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGHUP)
		<-c
		os.Exit(0)
	}()

	go debugserver.Start()

	err := dbconn.ConnectToDB("")
	if err != nil {
		log.Fatalf("Fatal error connecting to Postgres DB: %s", err)
	}

	protectedRoutes := http.NewServeMux()
	protectedRoutes.HandleFunc("/api/get", serveGet)
	protectedRoutes.HandleFunc("/api/update", serveUpdate)

	unprotectedRoutes := http.NewServeMux()
	unprotectedRoutes.Handle("/", http.FileServer(assets.Assets))
	unprotectedRoutes.Handle("/api/", AuthMiddleware(protectedRoutes))

	host := ""
	if env.InsecureDev {
		host = "127.0.0.1"
	}
	addr := net.JoinHostPort(host, port)
	log15.Info("management-console: listening", "addr", addr)

	s := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	unsafeNoHTTPS, _ := strconv.ParseBool(unsafeNoHTTPS)

	if unsafeNoHTTPS {
		s.Handler = unprotectedRoutes
		log.Fatalf("Fatal error serving: %s", s.ListenAndServe())
	}


	if err := configureTLS(); err != nil {
		log.Fatal("failed to configure TLS: error:", err)
	}
	s.Handler = HSTSMiddleware(unprotectedRoutes)
	log.Fatalf("Fatal error serving: %s", s.ListenAndServeTLS(tlsCert, tlsKey))
}

type jsonConfiguration struct {
	ID       string
	Contents string
}

func serveGet(w http.ResponseWriter, r *http.Request) {
	logger := log15.New("route", "get")

	critical, err := confdb.CriticalGetLatest(r.Context())
	if err != nil {
		logger.Error("confdb.CriticalGetLatest failed", "error", err)
		httpError(w, "Error retrieving latest configuration", "internal_error")
		return
	}

	err = json.NewEncoder(w).Encode(&jsonConfiguration{
		ID:       strconv.Itoa(int(critical.ID)),
		Contents: critical.Contents,
	})
	if err != nil {
		logger.Error("json response encoding failed", "error", err)
		httpError(w, "Error encoding json response", "internal_error")
	}
}

func httpError(w http.ResponseWriter, message string, code string) {
	_ = json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}{
		Error: message,
		Code:  code,
	})
}

func serveUpdate(w http.ResponseWriter, r *http.Request) {
	logger := log15.New("route", "update")

	var args struct {
		LastID   string
		Contents string
	}

	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		logger.Error("json argument decoding failed", "error", err)
		httpError(w, errors.Wrap(err, "Unexpected error when decoding arguments").Error(), "bad_request")
		return
	}

	if err != nil {

	}
}

func HSTSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// successfully authenticated.
		h.ServeHTTP(w, r)
	})
}
