package discovery

func (a *ARPDiscoverer) results() []ARPReply {
	a.mu.Lock()
	defer a.mu.Unlock()

	res := make([]ARPReply, 0, len(a.alive))
	for _, r := range a.alive {
		res = append(res, r)
	}
	return res
}
