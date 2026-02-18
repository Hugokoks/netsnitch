package tcp

import (
	"encoding/json"
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/output"
)

type TCPFormatter struct{}

func (f TCPFormatter) Protocol() domain.Protocol {
	return domain.TCP
}

func (f TCPFormatter) FormatRows(res domain.Result) string {

	var output string

	status := "CLOSED"
	if res.Open {
		status = "OPEN"
	}
	if res.Filtred {
		status = "FILTRED"

	}

	output = fmt.Sprintf(
		"[%s] %s:%d (%s) [%s]\n",
		status,
		res.IP,
		res.Port,
		res.Protocol,
		res.Service,
	)

	if res.Banner != "" {
		output += fmt.Sprintf(" └ banner: %s\n", res.Banner)
	}

	return output
}

func (f TCPFormatter) FormatJson(res domain.Result) string {

	output, _ := json.Marshal(res)
	return string(output) + "\n"

}

// Automatically register formatter on package load
func init() {
	output.Register(TCPFormatter{})
}
