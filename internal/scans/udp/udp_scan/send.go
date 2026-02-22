package udp_scan

import (
	"encoding/binary"
	"math/rand"
	"net"
	"netsnitch/internal/netutils"
	"netsnitch/internal/packet"
	"syscall"
)

func (m *Manager) sendUDP(dstIP net.IP, dstPort int) error {

	srcIP, err := netutils.GetLocalIP(dstIP)
	if err != nil {
		return err
	}

	srcPort := uint16(40000 + rand.Intn(20000))

	payload := []byte{}

	udpLen := 8 + len(payload)
	udp := make([]byte, udpLen)

	// Source port
	binary.BigEndian.PutUint16(udp[0:2], srcPort)

	// Destination port
	binary.BigEndian.PutUint16(udp[2:4], uint16(dstPort))

	// Length
	binary.BigEndian.PutUint16(udp[4:6], uint16(udpLen))

	// Checksum (zatím 0)
	binary.BigEndian.PutUint16(udp[6:8], 0)

	copy(udp[8:], payload)

	// Compute checksum
	cs := packet.TransportChecksum(srcIP, dstIP, udp, 17)
	binary.BigEndian.PutUint16(udp[6:8], cs)

	addr := &syscall.SockaddrInet4{
		Port: dstPort,
	}
	copy(addr.Addr[:], dstIP.To4())

	return syscall.Sendto(m.fdUDP, udp, 0, addr)
}
