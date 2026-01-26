package discovery

import (
	"fmt"
)

func (a *ARPDiscoverer) collect() {
	defer a.wg.Done()
	for reply := range a.replyChan {
		a.mu.Lock()
		if _, exists := a.alive[reply.IP.String()]; !exists {
			a.alive[reply.IP.String()] = reply.IP
			fmt.Printf("[ARP] ✓ %s (%s)\n", reply.IP, reply.MAC)
		}
		a.mu.Unlock()
	}
}
