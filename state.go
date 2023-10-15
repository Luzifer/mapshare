package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func loadState() error {
	if cfg.StateFile == "" {
		// No state file, no retaining
		return nil
	}

	if _, err := os.Stat(cfg.StateFile); err != nil {
		logrus.WithError(err).Warn("Unable to load state, using empty state")
		if os.IsNotExist(err) {
			return nil
		}
		return errors.Wrap(err, "accessing state file")
	}

	reqRetainerLock.Lock()
	defer reqRetainerLock.Unlock()

	f, err := os.Open(cfg.StateFile)
	if err != nil {
		return errors.Wrap(err, "opening state file")
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.WithError(err).Error("closing state file (leaked fd)")
		}
	}()

	return errors.Wrap(json.NewDecoder(f).Decode(&reqRetainer), "decoding state file")
}

func retainState() error {
	if cfg.StateFile == "" {
		// No state file, no retaining
		return nil
	}

	f, err := os.Create(cfg.StateFile)
	if err != nil {
		return errors.Wrap(err, "creating state file")
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.WithError(err).Error("closing state file (leaked fd)")
		}
	}()

	reqRetainerLock.RLock()
	defer reqRetainerLock.RUnlock()

	tmpState := make(map[string]position)
	for m, p := range reqRetainer {
		if time.Since(p.Time) > cfg.StateTimeout {
			continue
		}
		tmpState[m] = p
	}

	return errors.Wrap(json.NewEncoder(f).Encode(tmpState), "encoding state file")
}
