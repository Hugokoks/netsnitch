package udp_scan

import (
	"context"
	"fmt"
	"net"
	"time"
)

func (m *Manager) Scan(
	ctx context.Context,
	ip net.IP,
	port int,
	timeout time.Duration,
) UDPState {

	m.wg.Add(1)
	defer m.wg.Done()

	m.startOnce.Do(func() {
		go m.listen(m.fdUDP, m.handleUDP)
		go m.listen(m.fdICMP, m.handleICMP)

		go func() {
			m.wg.Wait()
			m.Close()
		}()
	})

	scanCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	key := fmt.Sprintf("%s:%d", ip.String(), port)
	respCh := make(chan UDPState, 1)

	// register pending
	m.mu.Lock()
	m.pending[key] = respCh
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()
	}()

	// send packet
	if err := m.sendUDP(ip, port); err != nil {
		return UDPUknown
	}

	select {

	case state := <-respCh:
		return state

	case <-scanCtx.Done():
		return UDPUknown // OPEN|FILTERED
	}
}
