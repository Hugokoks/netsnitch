package fingerprint

import (
	"bytes"
	"strings"
)

func (e *Engine) Detect(port int, raw string) *ServiceInfo {
	if raw == "" {
		return nil
	}

	rawLower := strings.ToLower(raw)
	rawBytes := []byte(raw)

	for _, r := range e.Rules {
		if len(r.Ports) > 0 {
			ok := false
			for _, p := range r.Ports {
				if p == port {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}

		/*TEST Performance and accuracy without When*/
		switch r.When.Type {
		case "prefix":
			if !strings.HasPrefix(raw, r.When.Pattern) {
				continue
			}

		case "contains":
			if !strings.Contains(rawLower, strings.ToLower(r.When.Pattern)) {
				continue
			}

		case "hex":
			if len(r.whenHex) == 0 {
				continue
			}
			if !bytes.Contains(rawBytes, r.whenHex) {
				continue
			}
		}

		if r.Match == nil || r.Match.Type == "" {
			return &ServiceInfo{
				Service:    r.Service,
				Product:    r.Product,
				Version:    "",
				Banner:     raw,
				Confidence: r.Confidence,
				RuleID:     r.ID,
			}
		}

		switch r.Match.Type {
		case "regex":

			if r.re == nil {
				continue
			}

			m := r.re.FindStringSubmatch(raw)
			if m == nil {
				continue
			}

			info := &ServiceInfo{
				Service:    r.Service,
				Product:    r.Product,
				Banner:     raw,
				Confidence: r.Confidence,
				RuleID:     r.ID,
			}

			if r.Extract.Version > 0 && r.Extract.Version < len(m) {
				info.Version = m[r.Extract.Version]
			}
			if r.Extract.Product > 0 && r.Extract.Product < len(m) {
				info.Product = m[r.Extract.Product]
			}

			return info

		default:
			continue
		}
	}

	return nil
}
