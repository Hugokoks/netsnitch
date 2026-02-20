package udp_scan

import (
	"fmt"
	"syscall"
)

func (m *Manager) listen(fd int, dispetcher func(packet []byte)) {

	buf := make([]byte, 65535)
	fmt.Println("listener start")

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
