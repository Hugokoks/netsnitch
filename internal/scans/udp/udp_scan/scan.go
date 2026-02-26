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
	payload []byte,
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
	respChan := make(chan UDPState, 1)

	m.mu.Lock()
	m.pending[key] = respChan
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()
	}()

	if err := m.sendUDP(ip, port, payload); err != nil {
		return UDPClosed
	}

	select {

	case state := <-respChan:

		return state

	case <-scanCtx.Done():
		// UDP timeout = open|filtered (unknown)
		return UDPOpenOrFiltered
	}
}
