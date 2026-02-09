package discovery

import (
	"context"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (a *ARPDiscoverer) sendARPRequest(ctx context.Context, handle *ARPHandle, ips []net.IP) {
	fmt.Println("[ARP] sending requests")

	for _, targetIP := range ips {
		select {
		case <-ctx.Done():
			return
		default:
			// Skip special IPs
			if isSpecialIP(targetIP, handle.srcIP) {
				continue
			}
			if err := writeARPRequest(handle, targetIP); err != nil {
				a.stats.Errors.Add(1)
			} else {
				a.stats.Sent.Add(1)
			}
		}
	}

}

// sendARPRequest sends an ARP request for the target IP
func writeARPRequest(handle *ARPHandle, targetIP net.IP) error {
	// Create buffer for serialization
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	// Build Ethernet layer
	eth := &layers.Ethernet{
		SrcMAC:       handle.srcMAC,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // Broadcast
		EthernetType: layers.EthernetTypeARP,
	}

	// Build ARP layer
	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   handle.srcMAC,
		SourceProtAddress: handle.srcIP.To4(),
		DstHwAddress:      []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Unknown
		DstProtAddress:    targetIP.To4(),
	}

	// Serialize layers
	if err := gopacket.SerializeLayers(buf, opts, eth, arp); err != nil {
		return fmt.Errorf("serialize packet: %w", err)
	}

	// Write packet to network
	if err := handle.handle.WritePacketData(buf.Bytes()); err != nil {
		return fmt.Errorf("write packet: %w", err)
	}

	return nil
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
