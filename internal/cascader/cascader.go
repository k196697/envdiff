// Package cascader merges multiple env maps in priority order,
// where earlier entries take precedence over later ones (cascade semantics).
package cascader

import "sort"

// Layer represents a named env map at a specific priority level.
type Layer struct {
	Name string
	Env  map[string]string
}

// Result holds the resolved value for a key across all layers.
type Result struct {
	Key        string
	Value      string
	SourceName string // name of the layer that provided the value
	Overridden []Override
}

// Override records a lower-priority layer that was shadowed.
type Override struct {
	LayerName string
	Value     string
}

// Cascade merges layers in order: layers[0] has the highest priority.
// For each key found in any layer, the first layer that defines it wins.
func Cascade(layers []Layer) []Result {
	if len(layers) == 0 {
		return nil
	}

	// Collect all unique keys across all layers.
	keySet := make(map[string]struct{})
	for _, l := range layers {
		for k := range l.Env {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]Result, 0, len(keys))
	for _, key := range keys {
		r := resolveKey(key, layers)
		results = append(results, r)
	}
	return results
}

func resolveKey(key string, layers []Layer) Result {
	r := Result{Key: key}
	for _, l := range layers {
		v, ok := l.Env[key]
		if !ok {
			continue
		}
		if r.SourceName == "" {
			// First (highest-priority) layer that has this key wins.
			r.Value = v
			r.SourceName = l.Name
		} else {
			r.Overridden = append(r.Overridden, Override{LayerName: l.Name, Value: v})
		}
	}
	return r
}
