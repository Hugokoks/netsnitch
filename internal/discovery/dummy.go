package discovery

import "net"

func openARPHandle(_ *net.Interface) (any, error) {
	return nil, nil
}

func sendARPRequest(_ any, _ *net.Interface, _ net.IP, _ net.IP) error {
	return nil
}

func readARPReply(_ any) net.IP {
	return nil
}

func mapToSlice(m map[string]net.IP) []net.IP {
	var ips []net.IP
	for _, ip := range m {
		ips = append(ips, ip)
	}
	return ips
}
