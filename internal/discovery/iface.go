package discovery

import (
	"errors"
	"net"
)

func pickInterface(ips []net.IP) (*net.Interface, net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, iface := range ifaces {
		// musí být up + mít MAC
		if iface.Flags & net.FlagUp == 0 || len(iface.HardwareAddr) == 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil || ip.To4() == nil {
				continue
			}

			// když interface IP sedí do stejné sítě jako target
			for _, targetIP := range ips {
				if sameSubnet(ip, targetIP) {
					return &iface, ip, nil
				}
			}
		}
	}

	return nil, nil, errors.New("no suitable interface found")
}

func sameSubnet(a, b net.IP) bool {
	a = a.To4()
	b = b.To4()
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}
