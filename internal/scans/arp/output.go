package arp_active

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/output"
)

type ARPFormatter struct{}

func (f ARPFormatter) Protocol() domain.Protocol {

	return domain.ARP
}

func (f ARPFormatter) Format(res domain.Result) {
	if !res.Alive {
		return
	}

	if res.MAC != nil {
		fmt.Printf("[ARP] %s (%s)\n", res.IP, res.MAC)
	} else {
		fmt.Printf("[ARP] %s\n", res.IP)
	}
}

func init() {
	output.Register(ARPFormatter{})

}
