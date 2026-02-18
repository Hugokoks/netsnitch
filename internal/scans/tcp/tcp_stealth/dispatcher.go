package tcp_stealth

import (
	"encoding/binary"
	"fmt"
	"net"
)

func (m *Manager) dispetcher(packet []byte) {

	if len(packet) < 20 {
		return
	}

	// extract IP headers bytes len from IP information byte
	// parse fist hlaf of bites
	ipHeaderLen := int(packet[0]&0x0F) * 4
	if ipHeaderLen < 20 || len(packet) < ipHeaderLen+20 {
		return
	}
	// extract tcp data based of ip header len
	tcp := packet[ipHeaderLen:]
	// 12-15 bytes always srcIP
	srcIP := net.IP(packet[12:16])
	srcPort := int(binary.BigEndian.Uint16(tcp[0:2]))

	flags := tcp[13]

	key := fmt.Sprintf("%s:%d", srcIP.String(), srcPort)
	ack := binary.BigEndian.Uint32(tcp[8:12])

	m.mu.Lock()
	conn, ok := m.pending[key]
	m.mu.Unlock()

	if !ok {
		return
	}
	// validate ACK number
	if ack != conn.seq+1 {
		return
	}

	if flags&0x12 == 0x12 { // 18 SYN+ACK
		conn.ch <- true
	} else if flags&0x04 == 0x04 { // 4 RST
		conn.ch <- false
	}
}
