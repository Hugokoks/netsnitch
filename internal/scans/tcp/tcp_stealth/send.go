package tcp_stealth

import (
	"encoding/binary"
	"math/rand"
	"net"
	"syscall"
)

//	----SYN PACKET ----
//
// [0]	Source Port (High)	0x9C	High byte of 40000 (0x9C40)
// [1]	Source Port (Low)	0x40	Low byte of 40000
// [2]	Dest Port (High)	0x00	High byte of 80 (0x0050)
// [3]	Dest Port (Low)	0x50	Low byte of 80
// [4]	Sequence No (B3)	0xDE	1st byte of random seq (e.g., 0xDEADBEEF)
// [5]	Sequence No (B2)	0xAD	2nd byte
// [6]	Sequence No (B1)	0xBE	3rd byte
// [7]	Sequence No (B0)	0xEF	4th byte
// [8]	Ack No (B3)	0x00	Always 0 in a SYN packet
// [9]	Ack No (B2)	0x00
// [10]	Ack No (B1)	0x00
// [11]	Ack No (B0)	0x00
// [12]	Data Offset	0x50	5 << 4 (Header is 20 bytes long)
// [13]	Flags	0x02	SYN flag enabled
// [14]	Window Size (High)	0xFF	High byte of 65535 (0xFFFF)
// [15]	Window Size (Low)	0xFF	Low byte of 65535
// [16]	Checksum (High)	0x??	Calculated by tcpChecksum function
// [17]	Checksum (Low)	0x??	Calculated by tcpChecksum function
// [18]	Urgent Ptr (High)	0x00	Not used (0)
// [19]	Urgent Ptr (Low)	0x00	Not used (0)

func (m *Manager) sendSYN(dstIP net.IP, dstPort int) error {

	// Get correct local IP used to reach destination
	srcIP, err := getLocalIP(dstIP)
	if err != nil {
		return err
	}

	// waiting for res 40000 - 60000
	srcPort := uint16(40000 + rand.Intn(20000))

	// res SYN/ACK seq + 1
	seq := rand.Uint32()

	tcp := make([]byte, 20)

	// Source port

	// Convert uint16 to 2 bytes
	binary.BigEndian.PutUint16(tcp[0:2], srcPort)

	// Destination port
	binary.BigEndian.PutUint16(tcp[2:4], uint16(dstPort))

	// Sequence number
	binary.BigEndian.PutUint32(tcp[4:8], seq)

	// Acknowledgment number (0 for SYN)
	// SYN/ACK Seq + 1
	binary.BigEndian.PutUint32(tcp[8:12], 0)

	// Data offset (5 * 4 = 20 bytes), no options
	// receiver knows where data starts
	tcp[12] = (5 << 4)

	// SYN flag 2
	tcp[13] = 0x02

	// Window size
	binary.BigEndian.PutUint16(tcp[14:16], 65535)

	// Checksum initially 0
	binary.BigEndian.PutUint16(tcp[16:18], 0)

	// Urgent pointer
	binary.BigEndian.PutUint16(tcp[18:20], 0)

	// Compute TCP checksum with pseudo-header
	cs := tcpChecksum(srcIP, dstIP, tcp)
	binary.BigEndian.PutUint16(tcp[16:18], cs)

	// Destination sockaddr
	addr := &syscall.SockaddrInet4{
		Port: dstPort,
	}
	copy(addr.Addr[:], dstIP.To4())

	// Send to network
	return syscall.Sendto(m.fd, tcp, 0, addr)
}

func getLocalIP(target net.IP) (net.IP, error) {
	// make fake UDP connection with target to obtain right ip of responsible network card
	conn, err := net.Dial("udp", target.String()+":80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
