// Package watcher provides a lightweight file-polling mechanism for .env files.
//
// It computes a SHA-256 hash of each watched file on every tick and fires a
// ChangeEvent callback when the hash differs from the previously recorded
// value. This makes it suitable for detecting edits, truncations, or full
// rewrites without relying on OS-level filesystem notifications.
//
// Example usage:
//
//	w := watcher.New([]string{".env", ".env.production"}, time.Second, func(e watcher.ChangeEvent) {
//		fmt.Printf("%s changed\n", e.Path)
//	})
//	w.Start()
//	defer w.Stop()
package watcher
