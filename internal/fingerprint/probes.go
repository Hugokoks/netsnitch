package fingerprint

type Probe struct {
	Name     string `json:"name"`
	Payload  string `json:"payload"`
	Ports    []int  `json:"ports"`
	Priority int    `json:"priority"`

	RawData []byte
}
