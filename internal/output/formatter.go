package output

import "netsnitch/internal/scan"

// Formatter defines how a scan result should be displayed.

type Formatter interface{


	Protocol() scan.Protocol
	Format(res scan.Result)

}


// internal registry of formatters by protocol

var formatters = make(map[scan.Protocol]Formatter)


// Register registers a formatter for a protocol.
// Called automatically via init() in formatter implementations.
func Register(f Formatter){

	formatters[f.Protocol()] = f

}