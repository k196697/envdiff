package watcher_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/user/envdiff/internal/watcher"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "KEY=original\n")

	var mu sync.Mutex
	var events []watcher.ChangeEvent

	w := watcher.New([]string{path}, 20*time.Millisecond, func(e watcher.ChangeEvent) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	})
	w.Start()
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(events) == 0 {
		t.Fatal("expected at least one ChangeEvent, got none")
	}
	if events[0].Path != path {
		t.Errorf("expected path %q, got %q", path, events[0].Path)
	}
	if events[0].OldHash == events[0].NewHash {
		t.Error("expected OldHash and NewHash to differ")
	}
}

func TestWatcher_NoChangeNoEvent(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "KEY=stable\n")

	var mu sync.Mutex
	var events []watcher.ChangeEvent

	w := watcher.New([]string{path}, 20*time.Millisecond, func(e watcher.ChangeEvent) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	})
	w.Start()
	defer w.Stop()

	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestWatcher_MissingFileIgnored(t *testing.T) {
	w := watcher.New([]string{"/nonexistent/.env"}, 20*time.Millisecond, func(e watcher.ChangeEvent) {
		t.Error("unexpected change event for missing file")
	})
	w.Start()
	time.Sleep(50 * time.Millisecond)
	w.Stop()
}
