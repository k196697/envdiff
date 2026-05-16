// Package validator checks .env values against simple type/format rules
// such as required, non-empty, numeric, URL, and boolean constraints.
package validator

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Rule describes a validation constraint for a specific key.
type Rule struct {
	Key      string
	Required bool
	Kind     string // "bool", "int", "float", "url", "nonempty"
}

// Violation records a failed rule for a key in a named file.
type Violation struct {
	File    string
	Key     string
	Rule    string
	Message string
}

// Validate applies rules to the provided diff results and returns any violations.
func Validate(results []diff.Result, rules []Rule) []Violation {
	if len(rules) == 0 || len(results) == 0 {
		return nil
	}

	ruleMap := make(map[string]Rule, len(rules))
	for _, r := range rules {
		ruleMap[strings.ToUpper(r.Key)] = r
	}

	var violations []Violation
	for _, res := range results {
		rule, ok := ruleMap[strings.ToUpper(res.Key)]
		if !ok {
			continue
		}
		for file, val := range res.Values {
			if rule.Required && val == "" {
				violations = append(violations, Violation{
					File:    file,
					Key:     res.Key,
					Rule:    "required",
					Message: "value is required but empty",
				})
				continue
			}
			if val == "" {
				continue
			}
			if v := checkKind(rule.Kind, val); v != "" {
				violations = append(violations, Violation{
					File:    file,
					Key:     res.Key,
					Rule:    rule.Kind,
					Message: fmt.Sprintf("%s: %s", rule.Kind, v),
				})
			}
		}
	}
	return violations
}

func checkKind(kind, val string) string {
	switch kind {
	case "bool":
		if _, err := strconv.ParseBool(val); err != nil {
			return fmt.Sprintf("%q is not a boolean", val)
		}
	case "int":
		if _, err := strconv.Atoi(val); err != nil {
			return fmt.Sprintf("%q is not an integer", val)
		}
	case "float":
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return fmt.Sprintf("%q is not a float", val)
		}
	case "url":
		u, err := url.ParseRequestURI(val)
		if err != nil || u.Scheme == "" {
			return fmt.Sprintf("%q is not a valid URL", val)
		}
	case "nonempty":
		if strings.TrimSpace(val) == "" {
			return "value must not be blank"
		}
	}
	return ""
}
