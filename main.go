package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Listen         string        `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		StateFile      string        `flag:"state-file" default:"" description:"Where to store retained locations (empty for no state)"`
		StateTimeout   time.Duration `flag:"state-timeout" default:"24h" description:"When to drop retained states"`
		VersionAndExit bool          `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing CLI options")
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log level")
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		fmt.Printf("mapshare %s\n", version) //nolint:forbidigo
		return
	}

	if err = loadState(); err != nil {
		logrus.WithError(err).Fatal("loading state")
	}

	r := mux.NewRouter()

	r.HandleFunc("/", handleRedirectRandom).Methods(http.MethodGet)

	r.PathPrefix("/asset/").Handler(
		http.StripPrefix("/asset/", http.FileServer(http.Dir("frontend"))),
	).Methods(http.MethodGet)

	r.HandleFunc("/{mapID}", handleMapFrontend).Methods(http.MethodGet)
	r.HandleFunc("/{mapID}", handleMapSubmit).Methods(http.MethodPut)
	r.HandleFunc("/{mapID}/ws", handleMapSocket).Methods(http.MethodGet)

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}

	logrus.WithField("version", version).Info("mapshare ready to serve")
	if err = server.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("running HTTP server")
	}
}
