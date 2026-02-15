package domain

import (
	"time"
)

type Config struct {
	Type    Protocol
	Timeout time.Duration

	Scope  Scope ////scope of ip's cidr, ip's,single ip
	Render RenderType
	Ports  PortScope
	Mode   ScanMode
}

func NewDefaultConfig() Config {

	return Config{
		Timeout: DefaultTimeout,
		Render:  ROWS_OUT,
	}

}

var DefaultPorts = []int{22, 80, 443}
var DefaultTimeout = 400 * time.Millisecond
