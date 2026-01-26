package discovery

import "context"

// listenARPReplies listens for ARP replies and sends them to the channel
func listenARPReplies(ctx context.Context, handle *ARPHandle, replies chan<- ARPReply, stats *Stats) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Read packet data (non-blocking with timeout)
			data, _, err := handle.handle.ReadPacketData()
			if err != nil {
				continue // Timeout or error, keep listening
			}

			// Parse ARP reply
			ip, mac, err := parseARPReply(data, handle.srcIP)
			if err != nil {
				continue // Not a valid reply
			}

			stats.Received.Add(1)

			// Send to channel (non-blocking)
			select {
			case replies <- ARPReply{IP: ip, MAC: mac}:
			case <-ctx.Done():
				return
			}
		}
	}
}
