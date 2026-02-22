package udp_scan

import (
	"context"
	"fmt"
	"net"
	"netsnitch/internal/domain"
	"time"
)

func (m *Manager) Scan(
	ctx context.Context,
	ip net.IP,
	port int,
	timeout time.Duration,
) domain.Result {

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

	// register pending
	m.mu.Lock()
	m.pending[key] = respChan
	m.mu.Unlock()

	// cleanup
	defer func() {
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()
	}()

	result := domain.Result{
		Protocol: domain.UDP,
		IP:       ip,
		Port:     port,
		Open:     false,
	}

	// send packet
	if err := m.sendUDP(ip, port); err != nil {
		return result
	}

	select {

	case state := <-respChan:
		switch state {
		case UDPOpen:
			result.Open = true
		case UDPClosed:
			result.Open = false
		}
		return result

	case <-scanCtx.Done():
		// timeout → OPEN|FILTERED
		result.Filtred = true
		return result
	}
}
