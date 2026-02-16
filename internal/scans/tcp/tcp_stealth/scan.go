package tcp_stealth

import (
	"fmt"
	"net"
	"netsnitch/internal/domain"
	"time"
)

func (m *Manager) Scan(
	ip net.IP,
	port int,
	timeout time.Duration,
) domain.Result {

	key := fmt.Sprintf("%s:%d", ip.String(), port)

	respCh := make(chan bool, 1)
	// --- register ---
	m.mu.Lock()
	m.pending[key] = respCh
	m.mu.Unlock()

	// --- send SYN ---
	err := m.sendSYN(ip, port)

	// cannot send tcp req to specific ip and port
	if err != nil {
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()

		return domain.Result{
			Protocol: domain.TCP,
			IP:       ip,
			Port:     port,
			Open:     false,
		}
	}

	// --- wait for response or timeout ---
	select {
	case open := <-respCh:
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()

		return domain.Result{
			Protocol: domain.TCP,
			IP:       ip,
			Port:     port,
			Open:     open,
		}

	case <-time.After(timeout):
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()

		return domain.Result{
			Protocol: domain.TCP,
			IP:       ip,
			Port:     port,
			Open:     false, // filtered / no response
		}
	}
}
