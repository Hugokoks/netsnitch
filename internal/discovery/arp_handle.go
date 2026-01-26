package discovery

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket/pcap"
)

// ARPHandle wraps a pcap handle for ARP packet capture and injection.
// It holds the raw pcap handle along with interface metadata needed
// for constructing and sending ARP requests.
type ARPHandle struct {
	handle *pcap.Handle     // Raw pcap handle for reading/writing packets
	iface  *net.Interface   // Network interface being used
	srcMAC net.HardwareAddr // Source MAC address (our MAC)
	srcIP  net.IP           // Source IP address (our IP)
}

// openARPHandle opens a pcap handle on the given interface with the specified source IP.
// The interface and source IP should be obtained from pickInterface() beforehand.
func openARPHandle(iface *net.Interface, srcIP net.IP) (*ARPHandle, error) {
	// Open pcap handle on interface
	handle, err := pcap.OpenLive(iface.Name, 65536, true, 10*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("pcap open failed (need root): %w", err)
	}

	// Set BPF filter to only capture ARP packets
	// This filtering happens in the kernel for efficiency
	if err := handle.SetBPFFilter("arp"); err != nil {
		handle.Close()
		return nil, fmt.Errorf("set BPF filter failed: %w", err)
	}

	return &ARPHandle{
		handle: handle,
		iface:  iface,
		srcMAC: iface.HardwareAddr,
		srcIP:  srcIP,
	}, nil
}

// Close closes the underlying pcap handle and releases resources.
func (h *ARPHandle) Close() error {
	h.handle.Close()
	return nil
}
