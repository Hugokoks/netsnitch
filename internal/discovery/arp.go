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
	Timeout time.Duration
	stats   Stats
}

func NewARP(timeout time.Duration) *ARPDiscoverer {
	return &ARPDiscoverer{Timeout: timeout}
}

// Discover performs ARP discovery on the given CIDR network
func (a *ARPDiscoverer) Discover(ctx context.Context, cidr string) ([]net.IP, error) {
	// 1. Parse CIDR
	ips, err := scan.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("parse CIDR: %w", err)
	}

	// 2. Pick interface + source IP
	iface, srcIP, err := pickInterface(ips)
	if err != nil {
		return nil, fmt.Errorf("pick interface: %w", err)
	}

	fmt.Printf("[ARP] Using interface %s (%s) with IP %s\n", iface.Name, iface.HardwareAddr, srcIP)

	// 3. Open ARP handle
	handle, err := openARPHandle(iface, srcIP)
	if err != nil {
		return nil, fmt.Errorf("open ARP handle: %w", err)
	}
	defer handle.Close()

	// 4. Setup channels and storage
	replyChan := make(chan ARPReply, 100)
	alive := make(map[string]net.IP)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create cancellable context with timeout
	scanCtx, cancel := context.WithTimeout(ctx, a.Timeout)
	defer cancel()

	// 5. Start listener goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		listenARPReplies(scanCtx, handle, replyChan, &a.stats)
	}()

	// 6. Start collector goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for reply := range replyChan {
			mu.Lock()
			if _, exists := alive[reply.IP.String()]; !exists {
				alive[reply.IP.String()] = reply.IP
				fmt.Printf("[ARP] ✓ %s (%s)\n", reply.IP, reply.MAC)
			}
			mu.Unlock()
		}
	}()

	// Wait a bit for listener to be ready
	time.Sleep(10 * time.Millisecond)

	// 7. Send ARP requests
	fmt.Printf("[ARP] Sending ARP requests to %d targets...\n", len(ips))
	for _, targetIP := range ips {
		select {
		case <-scanCtx.Done():
			goto cleanup
		default:
			// Skip special IPs
			if isSpecialIP(targetIP, srcIP) {
				continue
			}

			if err := sendARPRequest(handle, targetIP); err != nil {
				a.stats.Errors.Add(1)
			} else {
				a.stats.Sent.Add(1)
			}
		}
	}

cleanup:
	// 8. Wait for timeout or cancellation
	<-scanCtx.Done()

	// 9. Close channel and wait for goroutines
	close(replyChan)
	wg.Wait()

	// 10. Print stats and return
	fmt.Printf("[ARP] Stats - Sent: %d, Received: %d, Errors: %d\n",
		a.stats.Sent.Load(), a.stats.Received.Load(), a.stats.Errors.Load())
	fmt.Printf("[ARP] Found %d alive hosts\n", len(alive))

	mu.Lock()
	defer mu.Unlock()

	return mapToSlice(alive), nil
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

func mapToSlice(m map[string]net.IP) []net.IP {
	ips := make([]net.IP, 0, len(m))
	for _, ip := range m {
		ips = append(ips, ip)
	}
	return ips
}
