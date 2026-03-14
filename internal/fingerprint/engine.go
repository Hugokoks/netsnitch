package fingerprint

import (
	"path/filepath"
)

type Engine struct {
	ProbesByPort  map[int][]Probe
	GenericProbes []Probe

	portRules    []*Rule
	genericRules []*Rule

	prefixIndex   map[string][]*Rule
	containsIndex map[string][]*Rule
	hexIndex      map[string][]*Rule
	regexRules    []*Rule
}

func NewEngine() *Engine {
	e := &Engine{
		ProbesByPort:  make(map[int][]Probe),
		GenericProbes: []Probe{},
		portRules:     []*Rule{},
		genericRules:  []*Rule{},
		prefixIndex:   make(map[string][]*Rule),
		containsIndex: make(map[string][]*Rule),
		hexIndex:      make(map[string][]*Rule),
		regexRules:    []*Rule{},
	}

	return e
}

func InitFPEngine() (*Engine, error) {
	rulesPath := filepath.Join("data", "rules.json")
	probesPath := filepath.Join("data", "probes.json")

	fp := NewEngine()

	if err := fp.LoadRules(rulesPath); err != nil {
		return nil, err
	}
	if err := fp.LoadProbes(probesPath); err != nil {
		return nil, err
	}
	return fp, nil
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
