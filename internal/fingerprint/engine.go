package fingerprint

type Engine struct {
	ProbesByPort  map[int][]Probe
	GenericProbes []Probe

	Rules []Rule

	//prefixRules map[string][]Rule
	//tokenRules  map[string][]Rule
	//magicRules  [][]Rule
}

func NewEngine() *Engine {
	e := &Engine{}

	return e
}

func (e *Engine) getProbes(port int) []Probe {
	pp := e.ProbesByPort[port]
	gp := e.GenericProbes
	if len(gp) > 5 {
		gp = gp[:5]
	}

	out := make([]Probe, 0, len(pp)+len(gp))
	out = append(out, pp...)
	out = append(out, gp...)
	return out
}
