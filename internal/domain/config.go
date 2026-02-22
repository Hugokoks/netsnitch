package domain

import (
	"time"
)

type Config struct {
	Type    Protocol
	Timeout time.Duration

	Scope      Scope
	Render     RenderType
	Ports      PortScope
	Mode       ScanMode
	OpenOnly   bool
	UDPPayload string
}

func NewDefaultConfig() Config {

	return Config{
		Timeout: DefaultTimeout,
		Render:  ROWS_OUT,
	}

}

var DefaultPorts = []int{22, 80, 443}
var DefaultTimeout = 400 * time.Millisecond
