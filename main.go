package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

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

func init() {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("mapshare %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	if err := loadState(); err != nil {
		log.WithError(err).Fatal("Unable to load state")
	}

	r := mux.NewRouter()

	r.HandleFunc("/", handleRedirectRandom).Methods(http.MethodGet)

	r.PathPrefix("/asset/").Handler(
		http.StripPrefix("/asset/", http.FileServer(http.Dir("frontend"))),
	).Methods(http.MethodGet)

	r.HandleFunc("/{mapID}", handleMapFrontend).Methods(http.MethodGet)
	r.HandleFunc("/{mapID}", handleMapSubmit).Methods(http.MethodPut)
	r.HandleFunc("/{mapID}/ws", handleMapSocket).Methods(http.MethodGet)

	log.WithError(http.ListenAndServe(cfg.Listen, r)).Error("HTTP server caused an error")
}
