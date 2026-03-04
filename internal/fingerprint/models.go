package fingerprint

import "regexp"

type Pattern struct {
	Regex   *regexp.Regexp
	Service string
	Params  []XMLParam
	key     string
}

type ServiceInfo struct {
	Name    string
	Version string
	Product string
	Vendor  string
	Raw     string
}
