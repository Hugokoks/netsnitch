package fingerprint

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func (e *Engine) LoadRules(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read rules file: %w", err)
	}

	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return fmt.Errorf("unmarshal rules json: %w", err)
	}

	// Basic sanity: keep only rules with service + id
	filtered := make([]Rule, 0, len(rules))
	for i := range rules {

		r := &rules[i]

		if r.ID == "" || r.Service == "" {
			continue
		}

		///regex precompile
		if r.Match != nil && r.Match.Type == "regex" {
			re, err := regexp.Compile(r.Match.Pattern)
			if err != nil {
				continue
			}
			r.re = re
		}
		/// hex precompile
		if r.When != nil && r.When.Type == "hex" {

			sig, err := hex.DecodeString(strings.TrimSpace(r.When.Pattern))
			if err != nil || len(sig) == 0 {
				continue
			}
			r.whenHex = sig
		}

		filtered = append(filtered, *r)
	}

	e.Rules = filtered
	return nil
}
