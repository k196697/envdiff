package promoter

import (
	"encoding/json"
	"fmt"
	"io"
)

func writeText(w io.Writer, suggestions []Suggestion) error {
	if len(suggestions) == 0 {
		_, err := fmt.Fprintln(w, "No promotion candidates found.")
		return err
	}

	_, err := fmt.Fprintf(w, "Promotion candidates (%d):\n", len(suggestions))
	if err != nil {
		return err
	}

	for _, s := range suggestions {
		_, err = fmt.Fprintf(w, "  %-30s  %s → %s  (%s)\n",
			s.Key, s.FromEnv, s.ToEnv, s.Reason)
		if err != nil {
			return err
		}
	}
	return nil
}

type jsonSuggestion struct {
	Key       string `json:"key"`
	FromEnv   string `json:"from_env"`
	FromValue string `json:"from_value"`
	ToEnv     string `json:"to_env"`
	Reason    string `json:"reason"`
}

func writeJSON(w io.Writer, suggestions []Suggestion) error {
	out := make([]jsonSuggestion, 0, len(suggestions))
	for _, s := range suggestions {
		out = append(out, jsonSuggestion{
			Key:       s.Key,
			FromEnv:   s.FromEnv,
			FromValue: s.FromValue,
			ToEnv:     s.ToEnv,
			Reason:    s.Reason,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
