package fingerprint

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"regexp"
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

	e.portRules = []*Rule{}
	e.genericRules = []*Rule{}

	for i := range tempRules {
		r := tempRules[i]

		if r.ID == "" || r.Service == "" {
			continue
		}

		// precompile (Hex, Regex...)
		if r.When.Type == "hex" && r.When.Pattern != "" {
			if h, err := hex.DecodeString(r.When.Pattern); err == nil {
				r.whenHex = h
			}
		}

		if r.Match != nil && r.Match.Type == "regex" && r.Match.Pattern != "" {
			if re, err := regexp.Compile(r.Match.Pattern); err == nil {
				r.re = re
			}
		}

		if len(r.Ports) > 0 {

			newRule := r
			e.portRules = append(e.portRules, &newRule)
		} else {
			newRule := r
			e.genericRules = append(e.genericRules, &newRule)
		}
	}

	return nil
}
