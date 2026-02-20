package udp

import (
	"encoding/json"
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/output"
)

type UDPFormatter struct{}

func (f UDPFormatter) Protocol() domain.Protocol {
	return domain.UDP
}

func (f UDPFormatter) FormatRows(res domain.Result) string {

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
	return output
}

func (f UDPFormatter) FormatJson(res domain.Result) string {

	output, _ := json.Marshal(res)
	return string(output) + "\n"

}

func init() {
	output.Register(UDPFormatter{})
}
