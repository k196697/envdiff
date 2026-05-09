// Package profiler analyzes .env files and produces a profile summary
// describing key distribution, value patterns, and environment health.
package profiler

import (
	"sort"
	"strings"
)

// KeyProfile holds analysis data for a single key across all environments.
type KeyProfile struct {
	Key          string
	PresentIn    []string
	MissingIn    []string
	UniqueValues int
	HasEmpty     bool
}

// Profile is the result of profiling a set of environments.
type Profile struct {
	TotalKeys    int
	TotalEnvs    int
	FullCoverage int // keys present in all envs
	Partial      int // keys missing in at least one env
	AlwaysEmpty  int // keys that are empty in every env they appear
	Keys         []KeyProfile
}

// Analyze builds a Profile from a map of env name -> key/value pairs.
func Analyze(envs map[string]map[string]string) Profile {
	if len(envs) == 0 {
		return Profile{}
	}

	// Collect all unique keys.
	keySet := map[string]struct{}{}
	for _, kv := range envs {
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}

	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	profiles := make([]KeyProfile, 0, len(keys))
	fullCoverage, partial, alwaysEmpty := 0, 0, 0

	for _, key := range keys {
		kp := KeyProfile{Key: key}
		valueSet := map[string]struct{}{}
		emptyCount := 0

		for _, name := range envNames {
			val, ok := envs[name][key]
			if ok {
				kp.PresentIn = append(kp.PresentIn, name)
				valueSet[val] = struct{}{}
				if strings.TrimSpace(val) == "" {
					emptyCount++
				}
			} else {
				kp.MissingIn = append(kp.MissingIn, name)
			}
		}

		kp.UniqueValues = len(valueSet)
		kp.HasEmpty = emptyCount > 0

		if len(kp.MissingIn) == 0 {
			fullCoverage++
		} else {
			partial++
		}
		if emptyCount == len(kp.PresentIn) {
			alwaysEmpty++
		}

		profiles = append(profiles, kp)
	}

	return Profile{
		TotalKeys:    len(keys),
		TotalEnvs:    len(envNames),
		FullCoverage: fullCoverage,
		Partial:      partial,
		AlwaysEmpty:  alwaysEmpty,
		Keys:         profiles,
	}
}
