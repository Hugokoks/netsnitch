package fingerprint

import (
	"encoding/json"
	"fmt"
	"os"
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
	for _, r := range rules {
		if r.ID == "" || r.Service == "" {
			continue
		}
		filtered = append(filtered, r)
	}

	e.Rules = filtered
	return nil
}
