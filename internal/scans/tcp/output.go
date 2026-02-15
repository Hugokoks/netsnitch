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

	if !res.Open {
		return ""
	}

	output := fmt.Sprintf(
		"[OPEN] %s:%d (%s) [%s]\n",
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

	data, _ := json.Marshal(res)
	return string(data) + "\n"
}

// Automatically register formatter on package load
func init() {
	output.Register(TCPFormatter{})
}
