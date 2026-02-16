package tcp_stealth

import "syscall"

func (m *Manager) listen() {

	buf := make([]byte, 65535)

	for {
		select {
		case <-m.closeCh:
			return
		default:
		}

		n, _, err := syscall.Recvfrom(m.fd, buf, 0)
		if err != nil {
			continue
		}

		if n <= 0 {
			continue
		}

		m.dispetcher(buf[:n])
	}
}
