package output

import (
	"fmt"
	"netsnitch/internal/scan"
)

type TCPFormatter struct{}

func (f TCPFormatter) Protocol() scan.Protocol {
	return scan.TCP
}

func (f TCPFormatter) Format(res scan.Result) {
	if !res.Open {
		return
	}

	fmt.Printf(
		"[OPEN] %s:%d (%s) [%s]\n",
		res.IP,
		res.Port,
		res.Protocol,
		res.Service,
	)

	if res.Banner != "" {
		fmt.Printf(" └ banner: %s\n", res.Banner)
	}
}

// Automatically register formatter on package load
func init() {
	Register(TCPFormatter{})
}