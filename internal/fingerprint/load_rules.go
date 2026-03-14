package fingerprint

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"regexp"
	"strings"
)

func (e *Engine) LoadRules(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var tempRules []Rule
	if err := json.Unmarshal(data, &tempRules); err != nil {
		return err
	}

	// reset
	e.portRules = nil
	e.genericRules = nil
	e.prefixIndex = make(map[string][]*Rule)
	e.containsIndex = make(map[string][]*Rule)
	e.hexIndex = make(map[string][]*Rule)
	e.regexRules = nil

	for i := range tempRules {
		r := tempRules[i]

		if r.ID == "" || r.Service == "" || r.When == nil {
			continue
		}

		// precompile WHEN
		if r.When.Type == "hex" && r.When.Pattern != "" {
			if h, err := hex.DecodeString(r.When.Pattern); err == nil {
				r.whenHex = h
			} else {
				continue
			}
		}

		if r.When.Type == "regex" && r.When.Pattern != "" {
			if re, err := regexp.Compile(r.When.Pattern); err == nil {
				r.whenRe = re
			} else {
				continue
			}
		}

		// precompile MATCH
		if r.Match != nil && r.Match.Type == "regex" && r.Match.Pattern != "" {
			if re, err := regexp.Compile(r.Match.Pattern); err == nil {
				r.re = re
			} else {
				continue
			}
		}

		newRule := r

		if len(newRule.Ports) > 0 {
			e.portRules = append(e.portRules, &newRule)
		} else {
			e.genericRules = append(e.genericRules, &newRule)
		}

		e.addToIndex(&newRule)
	}

	return nil
}

func (e *Engine) addToIndex(r *Rule) {
	if r == nil || r.When == nil {
		return
	}

	switch r.When.Type {
	case "prefix":
		e.prefixIndex[r.When.Pattern] = append(e.prefixIndex[r.When.Pattern], r)

	case "contains":
		key := strings.ToLower(r.When.Pattern)
		e.containsIndex[key] = append(e.containsIndex[key], r)

	case "hex":
		e.hexIndex[r.When.Pattern] = append(e.hexIndex[r.When.Pattern], r)

	case "regex":
		e.regexRules = append(e.regexRules, r)
	}
}
