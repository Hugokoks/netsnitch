package output

import (
	"netsnitch/internal/domain"
)

// Formatter defines how a scan result should be displayed.

type Formatter interface{


	Protocol() domain.Protocol
	Format(res domain.Result)

}


// internal registry of formatters by protocol

var formatters = make(map[domain.Protocol]Formatter)


// Register registers a formatter for a protocol.
// Called automatically via init() in formatter implementations.
func Register(f Formatter){

	formatters[f.Protocol()] = f

}