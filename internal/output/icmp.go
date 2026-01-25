package output

import (
	"fmt"
	"netsnitch/internal/scan"
)

type ICMPFormatter struct{}

func (f ICMPFormatter) Protocol() scan.Protocol {
	return scan.ICMP
}

func (f ICMPFormatter) Format(res scan.Result) {
	if !res.Alive {
		return
	}

	fmt.Printf(
		"[HOST UP] %s  rtt=%v\n",
		res.IP,
		res.RTT,
	)
}

// Automatically register formatter on package load

func init(){

	Register(ICMPFormatter{})
}