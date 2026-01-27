package discovery

import "net"

func (a *ARPDiscoverer) results() []net.IP {
	a.mu.Lock()
	defer a.mu.Unlock()

	res := make([]net.IP, 0, len(a.alive))
	for _, ip := range a.alive {
		res = append(res, ip)
	}
	return res
}
