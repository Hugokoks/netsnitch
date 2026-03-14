package fingerprint

import (
	"bytes"
	"fmt"
	"strings"
)

func (e *Engine) Detect(port int, raw string) *ServiceInfo {
	if raw == "" {
		return nil
	}

	var best *ServiceInfo
	bestScore := -1

	candidates := e.candidateRules(port, raw)
	fmt.Printf("[detect] port=%d candidates=%d raw=%q\n", port, len(candidates), raw)
	for _, r := range candidates {
		fmt.Printf("[detect] try rule=%s service=%s when=%s pattern=%q\n",
			r.ID, r.Service, r.When.Type, r.When.Pattern)
		info := e.checkMatch(r, raw)
		if info == nil {
			continue
		}

		score := scoreMatch(info, containsPort(r.Ports, port))
		if score > bestScore {
			best = info
			bestScore = score
		}
	}
	if best == nil {
		fmt.Printf("[detect] miss port=%d raw=%q\n", port, raw)
	}

	return best
}
func (e *Engine) checkMatch(r *Rule, raw string) *ServiceInfo {
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

	case "regex":
		if r.whenRe != nil {
			matched = r.whenRe.MatchString(raw)
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

func (e *Engine) candidateRules(port int, raw string) []*Rule {
	var out []*Rule
	seen := make(map[*Rule]struct{})

	rawLower := strings.ToLower(raw)
	rawBytes := []byte(raw)

	addRules := func(rules []*Rule) {
		for _, r := range rules {
			if r == nil {
				continue
			}

			// port filter
			if len(r.Ports) > 0 && !containsPort(r.Ports, port) {
				continue
			}

			if _, ok := seen[r]; ok {
				continue
			}

			seen[r] = struct{}{}
			out = append(out, r)
		}
	}

	// prefix indexed rules
	for prefix, rules := range e.prefixIndex {
		if strings.HasPrefix(raw, prefix) {
			addRules(rules)
		}
	}

	// contains indexed rules
	for needle, rules := range e.containsIndex {
		if strings.Contains(rawLower, needle) {
			addRules(rules)
		}
	}

	// hex indexed rules
	for _, rules := range e.hexIndex {
		if len(rules) == 0 {
			continue
		}

		// všechny rules pod stejným klíčem mají stejný pattern,
		// takže stačí checknout první
		first := rules[0]
		if first == nil || len(first.whenHex) == 0 {
			continue
		}

		if bytes.Contains(rawBytes, first.whenHex) {
			addRules(rules)
		}
	}

	// regex when rules nejdou dobře indexovat, tak jako fallback
	addRules(e.regexRules)

	return out
}
