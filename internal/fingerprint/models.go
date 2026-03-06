package fingerprint

type ServiceInfo struct {
	Service    string  `json:"service"`
	Product    string  `json:"product,omitempty"`
	Version    string  `json:"version,omitempty"`
	Banner     string  `json:"banner,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	RuleID     string  `json:"rule_id,omitempty"`
}

type Rule struct {
	ID         string  `json:"id"`
	Service    string  `json:"service"`
	Product    string  `json:"product,omitempty"`
	Ports      []int   `json:"ports,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`

	When    *When   `json:"when,omitempty"`
	Match   *Match  `json:"match,omitempty"`
	Extract Extract `json:"extract,omitempty"`
}

type When struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

type Match struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
	Flags   string `json:"flags,omitempty"`
}

type Extract struct {
	Version int `json:"version,omitempty"`
	Product int `json:"product,omitempty"`
	Vendor  int `json:"vendor,omitempty"`
}

type Probe struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	Ports      []int  `json:"ports,omitempty"`
	Weight     int    `json:"weight,omitempty"`
	PayloadB64 string `json:"payload_b64,omitempty"`
	Payload    string `json:"payload,omitempty"`

	Raw []byte `json:"-"`
}
