package renamer

import (
	"encoding/json"
	"fmt"
	"io"
)

// Format controls the output format of WriteReport.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// WriteReport writes rename suggestions to w in the requested format.
func WriteReport(w io.Writer, suggestions []Suggestion, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, suggestions)
	default:
		return writeText(w, suggestions)
	}
}

func writeText(w io.Writer, suggestions []Suggestion) error {
	if len(suggestions) == 0 {
		_, err := fmt.Fprintln(w, "No rename suggestions found.")
		return err
	}
	for _, s := range suggestions {
		_, err := fmt.Fprintf(w, "  %-30s  →  %-30s  (score: %.2f)\n", s.From, s.To, s.Score)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, suggestions []Suggestion) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	type jsonSuggestion struct {
		From  string  `json:"from"`
		To    string  `json:"to"`
		Score float64 `json:"score"`
	}
	out := make([]jsonSuggestion, len(suggestions))
	for i, s := range suggestions {
		out[i] = jsonSuggestion{From: s.From, To: s.To, Score: s.Score}
	}
	return enc.Encode(out)
}
