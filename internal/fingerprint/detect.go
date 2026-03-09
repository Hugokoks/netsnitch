package fingerprint

import (
	"bytes"
	"strings"
)

func (e *Engine) Detect(port int, raw string) *ServiceInfo {
	if raw == "" {
		return nil
	}

	// Port rules
	for _, r := range e.portRules {
		if containsPort(r.Ports, port) {
			if info := e.checkMatch(r, raw); info != nil {
				return info
			}
		}
	}

	// Generic rules without ports
	for _, r := range e.genericRules {
		if info := e.checkMatch(r, raw); info != nil {
			return info
		}
	}

	return nil
}

func (e *Engine) checkMatch(r *Rule, raw string) *ServiceInfo {
	rawLower := strings.ToLower(raw)
	rawBytes := []byte(raw)
	matched := false

	switch r.When.Type {
	case "prefix":
		matched = strings.HasPrefix(raw, r.When.Pattern)
	case "contains":
		matched = strings.Contains(rawLower, strings.ToLower(r.When.Pattern))
	case "hex":
		if len(r.whenHex) > 0 {
			matched = bytes.Contains(rawBytes, r.whenHex)
		}
	}

	if !matched {
		return nil
	}

	info := &ServiceInfo{
		Service:    r.Service,
		Product:    r.Product,
		Banner:     raw,
		Confidence: r.Confidence,
		RuleID:     r.ID,
	}
	if r.Match == nil || r.Match.Type == "" {
		return info
	}

	switch r.Match.Type {

	case "regex":
		m := r.re.FindStringSubmatch(raw)
		if m != nil {
			if r.Extract.Version > 0 && r.Extract.Version < len(m) {
				info.Version = m[r.Extract.Version]
			}
			if r.Extract.Product > 0 && r.Extract.Product < len(m) {
				info.Product = m[r.Extract.Product]
			}
			return info
		}
	}
	return nil
}

// Pomocná funkce pro kontrolu portu v slice
func containsPort(ports []int, port int) bool {
	for _, p := range ports {
		if p == port {
			return true
		}
	}
	return false
}
