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
	var out string

	status := "CLOSED"
	if res.Open {
		status = "OPEN"
	}
	if res.Filtred {
		status = "FILTRED"
	}

	service := res.Service
	if service == "" {
		service = "unknown"
	}

	out = fmt.Sprintf(
		"[%s] %s:%d (%s) [%s]\n",
		status,
		res.IP,
		res.Port,
		res.Protocol,
		service,
	)

	if res.Product != "" {
		out += fmt.Sprintf(" └ product: %s\n", res.Product)
	}

	if res.Version != "" {
		out += fmt.Sprintf(" └ version: %s\n", res.Version)
	}

	if res.Banner != "" {
		label, value := output.FormatTextOrHex(res.Banner, 400, 64)
		out += fmt.Sprintf(" └ %s: %s\n", label, value)
	}

	return out
}

func (f TCPFormatter) FormatJson(res domain.Result) string {
	output, _ := json.Marshal(res)
	return string(output) + "\n"
}

func init() {
	output.Register(TCPFormatter{})
}
