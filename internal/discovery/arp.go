package discovery

import (
	"context"
	"fmt"
	"net"
	"netsnitch/internal/scan"
	"sync"
	"sync/atomic"
	"time"
)

// ARPReply represents a parsed ARP reply
type ARPReply struct {
	IP  net.IP
	MAC net.HardwareAddr
}

// Stats tracks ARP scanning statistics
type Stats struct {
	Sent     atomic.Uint64
	Received atomic.Uint64
	Errors   atomic.Uint64
}

type ARPDiscoverer struct {
	Timeout   time.Duration
	stats     Stats
	alive     map[string]net.IP
	mu        sync.Mutex
	replyChan chan ARPReply
	wg        sync.WaitGroup
}

func NewARP(timeout time.Duration) *ARPDiscoverer {
	return &ARPDiscoverer{
		Timeout:   timeout,
		alive:     make(map[string]net.IP),
		replyChan: make(chan ARPReply, 100),
	}
}

// Discover performs ARP discovery on the given CIDR network
func (a *ARPDiscoverer) Discover(ctx context.Context, cidr string) ([]net.IP, error) {
	//Parse CIDR
	ips, err := scan.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("parse CIDR: %w", err)
	}

	//Pick interface + source IP
	iface, srcIP, err := pickInterface(ips)
	if err != nil {
		return nil, fmt.Errorf("pick interface: %w", err)
	}

	fmt.Printf("[ARP] Using interface %s (%s) with IP %s\n", iface.Name, iface.HardwareAddr, srcIP)

	//Open ARP handle
	handle, err := openARPHandle(iface, srcIP)
	if err != nil {
		return nil, fmt.Errorf("open ARP handle: %w", err)
	}
	defer handle.Close()

	// Create cancellable context with timeout
	scanCtx, cancel := context.WithTimeout(ctx, a.Timeout)
	defer cancel()

	// Start listener goroutine
	a.wg.Add(1)
	go a.listenARPReplies(scanCtx, handle)

	//Start collector goroutine
	a.wg.Add(1)
	go a.collect()

	// Wait a bit for listener to be ready
	time.Sleep(10 * time.Millisecond)

	//Send ARP requests
	fmt.Printf("[ARP] Sending ARP requests to %d targets...\n", len(ips))
	a.sendRequest(ctx, handle, ips)

	//Wait for timeout or cancellation
	<-scanCtx.Done()

	//Close channel and wait for goroutines
	close(a.replyChan)
	a.wg.Wait()

	//Print stats and return
	fmt.Printf("[ARP] Stats - Sent: %d, Received: %d, Errors: %d\n",
		a.stats.Sent.Load(), a.stats.Received.Load(), a.stats.Errors.Load())
	fmt.Printf("[ARP] Found %d alive hosts\n", len(a.alive))

	a.mu.Lock()
	defer a.mu.Unlock()

	return a.results(), nil
}

// isSpecialIP checks if IP should be skipped (network, broadcast, own IP)
func isSpecialIP(ip, ourIP net.IP) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return true
	}

	// Network address (.0)
	if ip4[3] == 0 {
		return true
	}

	// Broadcast address (.255)
	if ip4[3] == 255 {
		return true
	}

	// Our own IP
	if ip.Equal(ourIP) {
		return true
	}

	return false
}

func (a *ARPDiscoverer) results() []net.IP {
	a.mu.Lock()
	defer a.mu.Unlock()

	res := make([]net.IP, 0, len(a.alive))
	for _, ip := range a.alive {
		res = append(res, ip)
	}
	return res
}
