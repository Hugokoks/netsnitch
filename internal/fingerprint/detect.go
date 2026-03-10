package fingerprint

import (
	"bytes"
	"strings"
)

func (e *Engine) Detect(port int, raw string) *ServiceInfo {
	if raw == "" {
		return nil
	}

	var best *ServiceInfo
	bestScore := -1

	// 1. Port-specific rules
	for _, r := range e.portRules {
		if !containsPort(r.Ports, port) {
			continue
		}

		info := e.checkMatch(r, raw)
		if info == nil {
			continue
		}

		score := scoreMatch(info, true)
		if score > bestScore {
			best = info
			bestScore = score
		}
	}

	// 2. Generic rules
	for _, r := range e.genericRules {
		info := e.checkMatch(r, raw)
		if info == nil {
			continue
		}

		score := scoreMatch(info, false)
		if score > bestScore {
			best = info
			bestScore = score
		}
	}

	return best
}
func (e *Engine) checkMatch(r *Rule, raw string) *ServiceInfo {
	///guards
	if r == nil || r.When == nil {
		return nil
	}

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
		if r.re == nil {
			return info
		}

		m := r.re.FindStringSubmatch(raw)
		if m == nil {
			return nil
		}

		if r.Extract.Version > 0 && r.Extract.Version < len(m) {
			info.Version = m[r.Extract.Version]
		}
		if r.Extract.Product > 0 && r.Extract.Product < len(m) {
			info.Product = m[r.Extract.Product]
		}

		return info
	}

	return nil
}

func containsPort(ports []int, port int) bool {
	for _, p := range ports {
		if p == port {
			return true
		}
	}
	return false
}

func scoreMatch(info *ServiceInfo, portSpecific bool) int {
	score := 0

	score += int(info.Confidence * 1000)

	// bonus for more informations
	if info.Product != "" {
		score += 100
	}
	if info.Version != "" {
		score += 100
	}

	// bonus for specific port
	if portSpecific {
		score += 50
	}

	return score
}
