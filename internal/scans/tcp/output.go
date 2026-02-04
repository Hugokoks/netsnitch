package tcp

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/output"
)

type TCPFormatter struct{}

func (f TCPFormatter) Protocol() domain.Protocol {
	return domain.TCP
}

func (f TCPFormatter) Format(res domain.Result) {

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
	output.Register(TCPFormatter{})
}
