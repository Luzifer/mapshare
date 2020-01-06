package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func loadState() error {
	if cfg.StateFile == "" {
		// No state file, no retaining
		return nil
	}

	if _, err := os.Stat(cfg.StateFile); err != nil {
		log.WithError(err).Warn("Unable to load state, using empty state")
		if os.IsNotExist(err) {
			return nil
		}
		return errors.Wrap(err, "Unable to access state file")
	}

	reqRetainerLock.Lock()
	defer reqRetainerLock.Unlock()

	f, err := os.Open(cfg.StateFile)
	if err != nil {
		return errors.Wrap(err, "Unable to open state file")
	}
	defer f.Close()

	return errors.Wrap(json.NewDecoder(f).Decode(&reqRetainer), "Unable to decode state file")

}

func retainState() error {
	if cfg.StateFile == "" {
		// No state file, no retaining
		return nil
	}

	f, err := os.Create(cfg.StateFile)
	if err != nil {
		return errors.Wrap(err, "Unable to create state file")
	}
	defer f.Close()

	reqRetainerLock.RLock()
	defer reqRetainerLock.RUnlock()

	var tmpState = make(map[string]position)
	for m, p := range reqRetainer {
		if time.Since(p.Time) > cfg.StateTimeout {
			continue
		}
		tmpState[m] = p
	}

	return errors.Wrap(json.NewEncoder(f).Encode(tmpState), "Unable to encode state file")
}
