package tcp_stealth

import (
	"fmt"
	"sync"
	"syscall"
)

type Manager struct {
	fd        socketFD ////file description
	pending   map[string]chan bool
	mu        sync.Mutex
	closeCh   chan struct{}
	wg        sync.WaitGroup
	startOnce sync.Once
}

func NewManager() (*Manager, error) {

	// Open raw TCP socket in kernel
	fd, err := syscall.Socket(
		syscall.AF_INET,     ////IPV4
		syscall.SOCK_RAW,    ////RAW packets
		syscall.IPPROTO_TCP, ////encapsulate IP headers on RAW TCP packet
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create raw socket: %w", err)
	}

	mgr := &Manager{
		fd:      fd,
		pending: make(map[string]chan bool),
		closeCh: make(chan struct{}),
	}

	return mgr, nil
}

func (m *Manager) Close() {
	close(m.closeCh)
	syscall.Close(m.fd)
}
