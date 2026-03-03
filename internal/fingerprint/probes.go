package fingerprint

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"sort"
)

type Probe struct {
	Name     string `json:"name"`
	Payload  string `json:"payload"` // Base64 encoded string from JSON
	Ports    []int  `json:"ports"`   // Preferred ports for this probe
	Priority int    `json:"priority"`
	RawData  []byte // Decoded binary data for the socket
}

// Global registry for O(1) port lookup
var probesByPort map[int][]Probe
var fallbackProbes []Probe // Probes that don't have a specific port

func LoadProbes(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var all []Probe
	if err := json.Unmarshal(file, &all); err != nil {
		return err
	}

	// Initialize our map
	probesByPort = make(map[int][]Probe)
	fallbackProbes = nil

	for i := range all {
		// 1. Pre-decode Base64 (as we discussed, for performance)
		data, _ := base64.StdEncoding.DecodeString(all[i].Payload)
		all[i].RawData = data

		// 2. Index by port
		if len(all[i].Ports) > 0 {
			for _, port := range all[i].Ports {
				probesByPort[port] = append(probesByPort[port], all[i])
			}
		} else {
			// Probes without specific ports (like GenericLine) go here
			fallbackProbes = append(fallbackProbes, all[i])
		}
	}

	// 3. Sort each port's slice by priority once at startup
	for port := range probesByPort {
		sort.Slice(probesByPort[port], func(i, j int) bool {
			return probesByPort[port][i].Priority < probesByPort[port][j].Priority
		})
	}

	return nil
}

func GetProbesForPort(port int) []Probe {
	// 1. Get port-specific probes from our pre-built hashmap
	// This is the O(1) lookup - extremely fast!
	prioritized := probesByPort[port]

	// 2. Create a combined list
	// We allocate space for both prioritized and fallback probes at once
	probes := make([]Probe, 0, len(prioritized)+len(fallbackProbes))

	// Add the "smart" guesses first
	probes = append(probes, prioritized...)

	// Add the "universal" keys last
	probes = append(probes, fallbackProbes...)

	return probes
}
