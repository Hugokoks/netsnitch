package tcp_stealth

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

	// active instance of task
	m.wg.Add(1)
	defer m.wg.Done()

	m.startOnce.Do(func() {
		// start listener only once for all tasks
		go m.listen()

		//close when there are no more active instances
		go func() {
			m.wg.Wait()
			m.Close()
		}()
	})

	scanCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	key := fmt.Sprintf("%s:%d", ip.String(), port)
	respCh := make(chan bool, 1)

	// --- register ---
	m.mu.Lock()
	m.pending[key] = respCh
	m.mu.Unlock()

	// clean pending map in end
	defer func() {
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()

	}()

	result := domain.Result{
		Protocol: domain.TCP,
		IP:       ip,
		Port:     port,
		Open:     false,
	}
	// --- send SYN ---
	err := m.sendSYN(ip, port)
	// cannot send tcp req to specific ip and port
	if err != nil {
		return result
	}
	// --- wait for response or timeout ---
	select {

	case open := <-respCh:
		result.Open = open
		return result

	case <-scanCtx.Done():
		return result
	}
}
