package input

import (
	"netsnitch/internal/domain"
)

type Query struct {
	Stages []Stage
}

type Stage struct {
	Protocol domain.Protocol
	Scope    domain.Scope
	Options  map[string]any
}
