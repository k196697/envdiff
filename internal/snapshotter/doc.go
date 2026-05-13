// Package snapshotter provides utilities for capturing, persisting, and
// comparing point-in-time snapshots of environment variable maps.
//
// A Snapshot records the full set of key/value pairs from one or more
// parsed .env files at a specific moment. Snapshots can be saved to disk
// as JSON and loaded back later for drift detection.
//
// Usage:
//
//	snap := snapshotter.Capture("production", envMap)
//	_ = snapshotter.Save(snap, ".envdiff/snapshots/prod.json")
//
//	prev, _ := snapshotter.Load(".envdiff/snapshots/prod.json")
//	deltas := snapshotter.Diff(prev, snap)
//	snapshotter.WriteReport(os.Stdout, prev, snap, deltas)
package snapshotter
