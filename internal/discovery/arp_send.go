package discovery

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// sendARPRequest sends an ARP request for the target IP
func sendARPRequest(handle *ARPHandle, targetIP net.IP) error {
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
