package udp_scan

import (
	"fmt"
	"netsnitch/internal/netutils"
	"syscall"
)

func (m *Manager) listen(fd netutils.SocketFD, dispetcher func(packet []byte)) {

	buf := make([]byte, 65535)

	for {
		select {
		case <-m.closeCh:
			return
		default:
		}

		n, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			dispetcher(buf[:n])
		}
	}
}
