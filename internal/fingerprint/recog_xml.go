package fingerprint

type RecogXML struct {
	Protocol     string           `xml:"protocol,attr"`
	Fingerprints []XMLFingerprint `xml:"fingerprint"`
}

type XMLFingerprint struct {
	Pattern     string     `xml:"pattern,attr"`
	Description string     `xml:"description"`
	Params      []XMLParam `xml:"param"`
}

type XMLParam struct {
	Name  string `xml:"name,attr"`
	Pos   int    `xml:"pos,attr"`
	Value string `xml:"value,attr"`
}
