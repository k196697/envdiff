package profiler_test

import (
	"testing"

	"github.com/user/envdiff/internal/profiler"
)

func baseEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"DB_HOST": "prod.db.local",
			"DB_PASS": "secret",
			"API_KEY": "abc123",
		},
		"staging": {
			"DB_HOST": "staging.db.local",
			"DB_PASS": "",
			"LOG_LEVEL": "debug",
		},
		"development": {
			"DB_HOST":   "localhost",
			"LOG_LEVEL": "debug",
		},
	}
}

func TestAnalyze_TotalsAreCorrect(t *testing.T) {
	p := profiler.Analyze(baseEnvs())

	if p.TotalEnvs != 3 {
		t.Errorf("TotalEnvs: got %d, want 3", p.TotalEnvs)
	}
	// Keys: DB_HOST, DB_PASS, API_KEY, LOG_LEVEL
	if p.TotalKeys != 4 {
		t.Errorf("TotalKeys: got %d, want 4", p.TotalKeys)
	}
}

func TestAnalyze_FullCoverageVsPartial(t *testing.T) {
	p := profiler.Analyze(baseEnvs())

	// DB_HOST is in all three → full coverage
	// DB_PASS, API_KEY, LOG_LEVEL are missing in at least one → partial
	if p.FullCoverage != 1 {
		t.Errorf("FullCoverage: got %d, want 1", p.FullCoverage)
	}
	if p.Partial != 3 {
		t.Errorf("Partial: got %d, want 3", p.Partial)
	}
}

func TestAnalyze_AlwaysEmpty(t *testing.T) {
	p := profiler.Analyze(baseEnvs())
	// DB_PASS is empty in staging and missing in development; only staging has it empty.
	// No key is empty in every env it appears across all three.
	if p.AlwaysEmpty != 0 {
		t.Errorf("AlwaysEmpty: got %d, want 0", p.AlwaysEmpty)
	}
}

func TestAnalyze_KeyProfileDetails(t *testing.T) {
	p := profiler.Analyze(baseEnvs())

	var dbHost *profiler.KeyProfile
	for i := range p.Keys {
		if p.Keys[i].Key == "DB_HOST" {
			dbHost = &p.Keys[i]
			break
		}
	}
	if dbHost == nil {
		t.Fatal("DB_HOST key profile not found")
	}
	if len(dbHost.PresentIn) != 3 {
		t.Errorf("DB_HOST PresentIn: got %d, want 3", len(dbHost.PresentIn))
	}
	if len(dbHost.MissingIn) != 0 {
		t.Errorf("DB_HOST MissingIn: got %d, want 0", len(dbHost.MissingIn))
	}
	// Three different values across envs.
	if dbHost.UniqueValues != 3 {
		t.Errorf("DB_HOST UniqueValues: got %d, want 3", dbHost.UniqueValues)
	}
}

func TestAnalyze_Empty(t *testing.T) {
	p := profiler.Analyze(map[string]map[string]string{})
	if p.TotalKeys != 0 || p.TotalEnvs != 0 {
		t.Errorf("expected empty profile, got %+v", p)
	}
}

func TestAnalyze_KeysSorted(t *testing.T) {
	p := profiler.Analyze(baseEnvs())
	for i := 1; i < len(p.Keys); i++ {
		if p.Keys[i].Key < p.Keys[i-1].Key {
			t.Errorf("keys not sorted: %s before %s", p.Keys[i-1].Key, p.Keys[i].Key)
		}
	}
}
