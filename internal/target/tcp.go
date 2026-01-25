package target

import (
	"net"
	"netsnitch/internal/scan"
)

type TCPBuilder struct{}

func (b TCPBuilder) Protocol() scan.Protocol{

	return scan.TCP
}


func (b TCPBuilder) Build(ips []net.IP) []Target{

		var targets []Target

		for _, ip := range ips {
			for _, port := range scan.DefaultPorts {
				targets = append(targets, Target{
					IP:   ip,
					Port: port,
				})
			}
		}

		return targets
}



func init(){

	Register(TCPBuilder{})
}