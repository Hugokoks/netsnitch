package arp_active

import (
	"encoding/json"
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/output"
)

type ARPFormatter struct{}

func (f ARPFormatter) Protocol() domain.Protocol {

	return domain.ARP
}

func (f ARPFormatter) FormatRows(res domain.Result) string {
	if !res.Alive {
		return ""
	}

	if res.MAC != nil {
		return fmt.Sprintf("[ARP] %s (%s)\n", res.IP, res.MAC)
	} else {
		return fmt.Sprintf("[ARP] %s\n", res.IP)
	}
}

func (f ARPFormatter) FormatJson(res domain.Result) string {
	tmp := struct {
		IP       string `json:"ip"`
		MAC      string `json:"mac,omitempty"`
		Protocol string `json:"protocol"`
		Alive    bool   `json:"is_alive"`
		RTT      int64  `json:"rtt_ns"`
	}{
		IP:       res.IP.String(),
		MAC:      res.MAC.String(),
		Protocol: "arp",
		Alive:    res.Alive,
		RTT:      res.RTT.Nanoseconds(),
	}
	data, _ := json.Marshal(tmp)

	return string(data) + "\n"
}

func init() {
	output.Register(ARPFormatter{})

}
