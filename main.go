package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Listen         string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func init() {
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
