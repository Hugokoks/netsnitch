package tcp_stealth

import (
	"encoding/binary"
	"fmt"
	"net"
)

func (m *Manager) dispetcher(packet []byte) {

	// extract IP headers bytes len from IP information byte
	// parse fist hlaf of bites
	ipHeaderLen := int(packet[0]&0x0F) * 4
	// extract tcp data based of ip header len
	tcp := packet[ipHeaderLen:]
	// 12-15 bytes always srcIP
	srcIP := net.IP(packet[12:16])
	srcPort := int(binary.BigEndian.Uint16(tcp[0:2]))

	flags := tcp[13]

	key := fmt.Sprintf("%s:%d", srcIP.String(), srcPort)

	m.mu.Lock()
	ch, ok := m.pending[key]
	m.mu.Unlock()

	if !ok {
		return
	}

	if flags&0x12 == 0x12 { // 18 SYN+ACK
		ch <- true
	} else if flags&0x04 == 0x04 { // 4 RST
		ch <- false
	}
}
