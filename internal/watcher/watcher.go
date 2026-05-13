// Package watcher monitors .env files for changes and triggers a callback
// when a modification is detected.
package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// ChangeEvent describes a file that has changed on disk.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
}

// Watcher polls a set of file paths and calls OnChange when any file's
// content hash differs from the last observed hash.
type Watcher struct {
	paths    []string
	interval time.Duration
	hashes   map[string]string
	OnChange func(ChangeEvent)
	stop     chan struct{}
}

// New creates a Watcher that checks the given paths every interval.
func New(paths []string, interval time.Duration, onChange func(ChangeEvent)) *Watcher {
	return &Watcher{
		paths:    paths,
		interval: interval,
		hashes:   make(map[string]string),
		OnChange: onChange,
		stop:     make(chan struct{}),
	}
}

// Start begins polling in a background goroutine. Call Stop to halt it.
func (w *Watcher) Start() {
	// Seed initial hashes so the first tick does not fire false positives.
	for _, p := range w.paths {
		if h, err := hashFile(p); err == nil {
			w.hashes[p] = h
		}
	}
	go w.loop()
}

// Stop signals the polling goroutine to exit.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) loop() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.check()
		case <-w.stop:
			return
		}
	}
}

func (w *Watcher) check() {
	for _, p := range w.paths {
		newHash, err := hashFile(p)
		if err != nil {
			continue
		}
		oldHash := w.hashes[p]
		if newHash != oldHash {
			w.hashes[p] = newHash
			if w.OnChange != nil {
				w.OnChange(ChangeEvent{Path: p, OldHash: oldHash, NewHash: newHash})
			}
		}
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
