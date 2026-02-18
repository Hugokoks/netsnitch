package discovery

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// parseARPReply parses raw packet data and extracts IP and MAC from ARP reply
func parseARPReply(data []byte, ourIP net.IP) (net.IP, net.HardwareAddr, error) {
	// Parse packet
	packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)

	// Get ARP layer
	arpLayer := packet.Layer(layers.LayerTypeARP)
	if arpLayer == nil {
		return nil, nil, fmt.Errorf("not an ARP packet")
	}

	arp := arpLayer.(*layers.ARP)

	// Check if it's a reply
	if arp.Operation != layers.ARPReply {
		return nil, nil, fmt.Errorf("not an ARP reply")
	}

	// verify reply is for our IP
	if !net.IP(arp.DstProtAddress).Equal(ourIP) {
		return nil, nil, fmt.Errorf("not for us")
	}

	// Extract sender IP and MAC
	senderIP := net.IP(arp.SourceProtAddress)
	senderMAC := net.HardwareAddr(arp.SourceHwAddress)

	// Skip our own IP
	if senderIP.Equal(ourIP) {
		return nil, nil, fmt.Errorf("reply from ourselves")
	}
	//fmt.Println("[ARP] reply from", senderIP, senderMAC)

	return senderIP, senderMAC, nil
}
