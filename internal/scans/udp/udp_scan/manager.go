package udp_scan

import (
	"fmt"
	"sync"
	"syscall"
)

type UDPState int

const (
	UDPUknown UDPState = iota
	UDPOpen
	UDPClosed
)

type Manager struct {
	fdUDP     int
	fdICMP    int
	mu        sync.Mutex
	closeCh   chan struct{}
	pending   map[string]chan UDPState
	wg        sync.WaitGroup
	startOnce sync.Once
}

func newManager() (*Manager, error) {

	fdUDP, err := syscall.Socket(
		syscall.AF_INET,     // IPV4
		syscall.SOCK_RAW,    // RAW SOCKET
		syscall.IPPROTO_UDP, // RAW UDP
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create raw socket: %w", err)
	}

	fdICMP, err := syscall.Socket(
		syscall.AF_INET,      // IPV4
		syscall.SOCK_RAW,     // RAW SOCKET
		syscall.IPPROTO_ICMP, // ICMP Protocol
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create raw socket: %w", err)
	}

	syscall.SetNonblock(fdUDP, true)
	syscall.SetNonblock(fdICMP, true)

	mgr := &Manager{
		fdUDP:   fdUDP,
		fdICMP:  fdICMP,
		closeCh: make(chan struct{}),
		pending: make(map[string]chan UDPState),
	}

	return mgr, nil

}

func (m *Manager) Close() {
	close(m.closeCh)
	syscall.Close(m.fdICMP)
	syscall.Close(m.fdUDP)

}
