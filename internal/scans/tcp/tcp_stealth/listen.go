package tcp_stealth

import (
	"fmt"
	"syscall"
)

func (m *Manager) listen() {

	buf := make([]byte, 65535)
	fmt.Println("listener start")
	for {
		select {
		case <-m.closeCh:
			return
		default:
		}

		n, _, err := syscall.Recvfrom(m.fd, buf, 0)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			m.dispetcher(buf[:n])

		}

		//
	}
}
