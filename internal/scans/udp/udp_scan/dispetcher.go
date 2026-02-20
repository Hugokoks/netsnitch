package udp_scan

import (
	"encoding/binary"
	"fmt"
	"net"
)

func (m *Manager) handleUDP(packet []byte) {

	payload, ok := ipPayload(packet)
	if !ok || len(payload) < 8 {
		return
	}

	// IP source je vždy na 12–15
	srcIP := net.IP(packet[12:16])
	srcPort := int(binary.BigEndian.Uint16(payload[0:2]))

	key := fmt.Sprintf("%s:%d", srcIP.String(), srcPort)

	m.mu.Lock()
	ch, exists := m.pending[key]
	m.mu.Unlock()

	if !exists {
		return
	}

	select {
	case ch <- UDPOpen:
	default:
	}
}

func (m *Manager) handleICMP(packet []byte) {

	payload, ok := ipPayload(packet)
	if !ok || len(payload) < 8 {
		return
	}

	icmpType := payload[0]
	icmpCode := payload[1]

	// Only Port Unreachable
	if icmpType != 3 || icmpCode != 3 {
		return
	}

	// Embedded original IP starts at payload[8:]
	if len(payload) < 8+20 {
		return
	}

	originalIP := payload[8:]

	origIPHeaderLen := int(originalIP[0]&0x0F) * 4
	if origIPHeaderLen < 20 || len(originalIP) < origIPHeaderLen+8 {
		return
	}

	originalUDP := originalIP[origIPHeaderLen:]

	dstPort := int(binary.BigEndian.Uint16(originalUDP[2:4]))
	srcIP := net.IP(originalIP[12:16])

	key := fmt.Sprintf("%s:%d", srcIP.String(), dstPort)

	m.mu.Lock()
	ch, exists := m.pending[key]
	m.mu.Unlock()

	if !exists {
		return
	}

	select {
	case ch <- UDPClosed:
	default:
	}
}

func ipPayload(packet []byte) ([]byte, bool) {
	if len(packet) < 20 {
		return nil, false
	}

	ipHeaderLen := int(packet[0]&0x0F) * 4
	if ipHeaderLen < 20 || len(packet) < ipHeaderLen {
		return nil, false
	}

	return packet[ipHeaderLen:], true
}
