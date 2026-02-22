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

	var status string

	if res.Open {
		status = "OPEN"
	} else if res.Service == "open|filtered" {
		status = "OPEN|FILTERED"
	} else {
		status = "CLOSED"
	}

	return fmt.Sprintf(
		"[%s] %s:%d (%s)\n",
		status,
		res.IP,
		res.Port,
		res.Protocol,
	)
}

func (f UDPFormatter) FormatJson(res domain.Result) string {

	output, _ := json.Marshal(res)
	return string(output) + "\n"

}

func init() {
	output.Register(UDPFormatter{})
}
