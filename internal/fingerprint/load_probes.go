package fingerprint

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func (e *Engine) LoadProbes(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read probes file: %w", err)
	}

	var probes []Probe
	if err := json.Unmarshal(data, &probes); err != nil {
		return fmt.Errorf("unmarshal probes json: %w", err)
	}

	// Reset
	e.ProbesByPort = make(map[int][]Probe)
	e.GenericProbes = nil

	for i := range probes {
		p := probes[i]

		if p.ID == "" {
			continue
		}

		// Default weight if not provided
		if p.Weight == 0 {
			p.Weight = 1
		}

		// Build Raw payload once
		raw, err := buildProbeRaw(p)
		if err != nil {
			// invalid base64 etc -> skip probe
			continue
		}
		p.Raw = raw

		// Store
		if len(p.Ports) == 0 {
			e.GenericProbes = append(e.GenericProbes, p)
			continue
		}

		for _, port := range p.Ports {
			e.ProbesByPort[port] = append(e.ProbesByPort[port], p)
		}
	}

	// Sort: higher Weight first (more important first)
	for port := range e.ProbesByPort {
		sort.SliceStable(e.ProbesByPort[port], func(i, j int) bool {
			return e.ProbesByPort[port][i].Weight > e.ProbesByPort[port][j].Weight
		})
	}
	sort.SliceStable(e.GenericProbes, func(i, j int) bool {
		return e.GenericProbes[i].Weight > e.GenericProbes[j].Weight
	})

	return nil
}

// buildProbeRaw prepares payload bytes for socket write.
// Supports payload_b64 (binary) and payload (text).
func buildProbeRaw(p Probe) ([]byte, error) {
	// Prefer base64 if present
	if strings.TrimSpace(p.PayloadB64) != "" {
		// Some payloads might include spaces/newlines; remove them.
		clean := strings.ReplaceAll(p.PayloadB64, " ", "")
		clean = strings.ReplaceAll(clean, "\n", "")
		clean = strings.ReplaceAll(clean, "\r", "")
		raw, err := base64.StdEncoding.DecodeString(clean)
		if err != nil {
			return nil, fmt.Errorf("probe %s invalid payload_b64: %w", p.ID, err)
		}
		return raw, nil
	}

	// Fallback to text payload
	if p.Payload != "" {
		return []byte(p.Payload), nil
	}

	// Empty payload is allowed (e.g., "null probe" / just read banner)
	return []byte{}, nil
}
