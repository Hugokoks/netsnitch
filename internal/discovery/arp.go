package discovery

import (
	"context"
	"fmt"
	"net"
	"netsnitch/internal/domain"
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

func (a *ARPDiscoverer) setup(cidr string) (*ARPHandle, []net.IP, error) {

	ips, err := domain.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, err
	}

	iface, srcIP, err :=  PickInterface(ips)
	if err != nil {
		return nil, nil, err
	}

	handle, err := openARPHandle(iface, srcIP)
	return handle, ips, err
}

// Discover performs ARP discovery on the given CIDR network
func (a *ARPDiscoverer) Discover(ctx context.Context, cidr string, cfg domain.Config) ([]net.IP, error) {

	handle, ips, err := a.setup(cidr)
	if err != nil {

		return nil, fmt.Errorf("arp setup error: %w", err)
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

	if cfg.Type == domain.ARP_ACTIVE {
		a.sendARPRequest(ctx, handle, ips)
	}

	//Wait for timeout or cancellation
	<-scanCtx.Done()

	a.stop()

	return a.results(), nil
}

func (a *ARPDiscoverer) stop() {

	//Close channel and wait for goroutines
	close(a.replyChan)
	a.wg.Wait()

	//Print stats and return
	fmt.Printf("[ARP] Stats - Sent: %d, Received: %d, Errors: %d\n",
		a.stats.Sent.Load(), a.stats.Received.Load(), a.stats.Errors.Load())
	fmt.Printf("[ARP] Found %d alive hosts\n", len(a.alive))

}
