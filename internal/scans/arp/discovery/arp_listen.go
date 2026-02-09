package discovery

import (
	"context"
	"fmt"
	"time"
)

// listenARPReplies listens for ARP replies and sends them to the channel
func (a *ARPDiscoverer) listenARPReplies(ctx context.Context, handle *ARPHandle) {
	defer a.wg.Done()
	fmt.Println("[ARP] listening...")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Read packet data (non-blocking with timeout)
			data, _, err := handle.handle.ReadPacketData()
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue // Timeout or error, keep listening
			}

			// Parse ARP reply
			ip, mac, err := parseARPReply(data, handle.srcIP)
			if err != nil {
				continue // Not a valid reply
			}

			a.stats.Received.Add(1)

			// Send to channel (non-blocking)
			select {
			case a.replyChan <- ARPReply{IP: ip, MAC: mac}:
			case <-ctx.Done():
				return
			}
		}
	}
}
