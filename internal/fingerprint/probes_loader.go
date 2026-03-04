package fingerprint

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"sort"
)

var probesByPort map[int][]Probe
var fallbackProbes []Probe

func LoadProbes(path string) error {

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var all []Probe

	err = json.Unmarshal(file, &all)
	if err != nil {
		return err
	}

	probesByPort = make(map[int][]Probe)

	for i := range all {

		data, _ := base64.StdEncoding.DecodeString(all[i].Payload)
		all[i].RawData = data

		if len(all[i].Ports) > 0 {

			for _, port := range all[i].Ports {

				probesByPort[port] =
					append(probesByPort[port], all[i])
			}

		} else {

			fallbackProbes =
				append(fallbackProbes, all[i])
		}
	}

	for port := range probesByPort {

		sort.Slice(probesByPort[port], func(i, j int) bool {

			return probesByPort[port][i].Priority <
				probesByPort[port][j].Priority
		})
	}

	return nil
}

func GetProbesForPort(port int) []Probe {

	prioritized := probesByPort[port]

	probes := make([]Probe, 0,
		len(prioritized)+len(fallbackProbes))

	probes = append(probes, prioritized...)
	probes = append(probes, fallbackProbes...)

	return probes
}
