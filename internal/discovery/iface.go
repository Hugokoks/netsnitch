package discovery

import (
	"errors"
	"net"
)

func PickInterface(ips []net.IP) (*net.Interface, net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, iface := range ifaces {

		// must be up + have MAC

		if iface.Flags&net.FlagUp == 0 || len(iface.HardwareAddr) == 0 {
			continue
		}
		addrs, _ := iface.Addrs()

		for _, addr := range addrs {
			ip, ipNet, err := net.ParseCIDR(addr.String())

			if err != nil || ip.To4() == nil {
				continue
			}

			// if the interface IP is in the same network as the target
			for _, targetIP := range ips {

				if ipNet.Contains(targetIP) {
					return &iface, ip, nil
				}
			}
		}
	}

	return nil, nil, errors.New("no suitable interface found")
}
