package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	websocketBufferSize     = 1024
	websocketUpdatePoolSize = 10
)

type position struct {
	Lat      float64   `json:"lat"`
	Lon      float64   `json:"lon"`
	Retained bool      `json:"retained"`
	SenderID string    `json:"sender_id"`
	Time     time.Time `json:"time"`
}

var (
	reqDistributors     = map[string]map[string]chan position{}
	reqDistributorsLock = new(sync.RWMutex)
	reqRetainer         = map[string]position{}
	reqRetainerLock     = new(sync.RWMutex)

	upgrader = websocket.Upgrader{
		ReadBufferSize:  websocketBufferSize,
		WriteBufferSize: websocketBufferSize,
	}
)

func handleRedirectRandom(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, uuid.Must(uuid.NewV4()).String(), http.StatusFound)
}

func handleMapFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/index.html")
}

func handleMapSocket(w http.ResponseWriter, r *http.Request) {
	var (
		vars    = mux.Vars(r)
		mapID   = vars["mapID"]
		sockID  = uuid.Must(uuid.NewV4()).String()
		updates = make(chan position, websocketUpdatePoolSize)
	)

	// Register update channel
	reqDistributorsLock.Lock()
	if _, ok := reqDistributors[mapID]; !ok {
		reqDistributors[mapID] = make(map[string]chan position)
	}
	reqDistributors[mapID][sockID] = updates
	reqDistributorsLock.Unlock()

	// In case a retained position is available queue it
	reqRetainerLock.RLock()
	if p, ok := reqRetainer[mapID]; ok && time.Since(p.Time) < cfg.StateTimeout {
		updates <- p
	}
	reqRetainerLock.RUnlock()

	// Queue deregistration
	defer func() {
		reqDistributorsLock.Lock()
		defer reqDistributorsLock.Unlock()

		delete(reqDistributors[mapID], sockID)
		close(updates)
	}()

	// Open socket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Debug("opening websocket")
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.WithError(err).Error("closing socket connection (leaked fd)")
		}
	}()

	// Send updates
	for pos := range updates {
		if err = conn.WriteJSON(pos); err != nil {
			logrus.WithError(err).Debug("sending position")
			return
		}
	}
}

func handleMapSubmit(w http.ResponseWriter, r *http.Request) {
	var (
		pos   position
		vars  = mux.Vars(r)
		mapID = vars["mapID"]
	)

	if err := json.NewDecoder(r.Body).Decode(&pos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pos.Time = time.Now()

	reqRetainerLock.Lock()
	if pos.Retained {
		reqRetainer[mapID] = pos
	} else {
		delete(reqRetainer, mapID)
	}
	reqRetainerLock.Unlock()

	go func() {
		if err := retainState(); err != nil {
			logrus.WithError(err).Error("retaining state to disk")
		}
	}()

	reqDistributorsLock.RLock()
	defer reqDistributorsLock.RUnlock()

	distributors, ok := reqDistributors[mapID]
	if !ok || len(distributors) == 0 {
		// No subscribers at all
		return
	}

	for _, c := range distributors {
		c <- pos
	}
}
