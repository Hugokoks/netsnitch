package fingerprint

import "strings"

// Detect runs the raw string through all loaded fingerprint patterns
func Detect(raw string) *ServiceInfo {
	if raw == "" {
		return nil
	}

	// Basic cleanup
	cleanRaw := strings.TrimSpace(raw)
	// Remove null bytes and other non-printable garbage common in banners
	cleanRaw = strings.Map(func(r rune) rune {
		if r < 32 || r > 126 {
			return -1
		}
		return r
	}, cleanRaw)

	if len(cleanRaw) < 2 {
		return nil
	}

	// Loop through our "brain"
	for _, p := range fingerprints {
		matches := p.Regex.FindStringSubmatch(cleanRaw)
		if matches == nil {
			continue
		}

		info := &ServiceInfo{
			Name: p.Service,
			Raw:  cleanRaw,
		}

		// If the regex has a capture group (e.g. ^SSH-(.+)), we take it as version
		if len(matches) > 1 {
			info.Version = matches[1]
		}

		return info
	}

	// If we got some data but no regex matched
	return &ServiceInfo{
		Name: "unknown",
		Raw:  cleanRaw,
	}
}
