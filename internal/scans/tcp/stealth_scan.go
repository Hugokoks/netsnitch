package tcp

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"netsnitch/internal/domain"
	"sync"
	"syscall"
	"time"
)

type StealthManager struct {
	fd      int ////file description
	pending map[string]chan bool
	mu      sync.Mutex
	closeCh chan struct{}
}

func NewStealthManager() (*StealthManager, error) {

	/////Open raw TCP socket
	fd, err := syscall.Socket(
		syscall.AF_INET,     ////IPV4
		syscall.SOCK_RAW,    ////Access to packets
		syscall.IPPROTO_TCP, ////TCP protocol
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create raw socket: %w", err)
	}

	mgr := &StealthManager{
		fd:      fd,
		pending: make(map[string]chan bool),
		closeCh: make(chan struct{}),
	}
	////Listen for answers
	go mgr.listen()

	return mgr, nil
}

func (m *StealthManager) listen() {

	buf := make([]byte, 65535)

	for {
		select {
		case <-m.closeCh:
			return
		default:
		}

		n, _, err := syscall.Recvfrom(m.fd, buf, 0)
		if err != nil {
			continue
		}

		if n <= 0 {
			continue
		}

		m.handlePacket(buf[:n])
	}
}
func (m *StealthManager) sendSYN(dstIP net.IP, dstPort int) error {

	// Get correct local IP used to reach destination
	srcIP, err := getLocalIP(dstIP)
	if err != nil {
		return err
	}

	srcPort := uint16(40000 + rand.Intn(20000))
	seq := rand.Uint32()

	tcp := make([]byte, 20)

	// Source port
	binary.BigEndian.PutUint16(tcp[0:2], srcPort)

	// Destination port
	binary.BigEndian.PutUint16(tcp[2:4], uint16(dstPort))

	// Sequence number
	binary.BigEndian.PutUint32(tcp[4:8], seq)

	// Acknowledgment number (0 for SYN)
	binary.BigEndian.PutUint32(tcp[8:12], 0)

	// Data offset (5 * 4 = 20 bytes), no options
	tcp[12] = (5 << 4)

	// SYN flag
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

	return syscall.Sendto(m.fd, tcp, 0, addr)
}

func (m *StealthManager) handlePacket(packet []byte) {

	ipHeaderLen := int(packet[0]&0x0F) * 4
	tcp := packet[ipHeaderLen:]

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

	if flags&0x12 == 0x12 { // SYN+ACK
		ch <- true
	} else if flags&0x04 == 0x04 { // RST
		ch <- false
	}
}
func (m *StealthManager) Scan(
	ip net.IP,
	port int,
	timeout time.Duration,
) domain.Result {

	key := fmt.Sprintf("%s:%d", ip.String(), port)

	respCh := make(chan bool, 1)

	// --- register ---
	m.mu.Lock()
	m.pending[key] = respCh
	m.mu.Unlock()

	// --- send SYN ---
	err := m.sendSYN(ip, port)

	////cannot send tcp req to specific ip and port
	if err != nil {
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()

		return domain.Result{
			Protocol: domain.TCP,
			IP:       ip,
			Port:     port,
			Open:     false,
		}
	}

	// --- wait for response or timeout ---
	select {
	case open := <-respCh:
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()

		return domain.Result{
			Protocol: domain.TCP,
			IP:       ip,
			Port:     port,
			Open:     open,
		}

	case <-time.After(timeout):
		m.mu.Lock()
		delete(m.pending, key)
		m.mu.Unlock()

		return domain.Result{
			Protocol: domain.TCP,
			IP:       ip,
			Port:     port,
			Open:     false, // filtered / no response
		}
	}
}
func tcpChecksum(srcIP, dstIP net.IP, tcp []byte) uint16 {

	pseudoHeader := make([]byte, 12+len(tcp))

	// source IP
	copy(pseudoHeader[0:4], srcIP.To4())

	// dest IP
	copy(pseudoHeader[4:8], dstIP.To4())

	// zero
	pseudoHeader[8] = 0

	// protocol (TCP = 6)
	pseudoHeader[9] = 6

	// TCP length
	binary.BigEndian.PutUint16(pseudoHeader[10:12], uint16(len(tcp)))

	// TCP header
	copy(pseudoHeader[12:], tcp)

	return checksum(pseudoHeader)
}

func checksum(data []byte) uint16 {
	var sum uint32

	for i := 0; i+1 < len(data); i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i:]))
	}

	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}

	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	return ^uint16(sum)
}

func getLocalIP(target net.IP) (net.IP, error) {
	conn, err := net.Dial("udp", target.String()+":80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func (m *StealthManager) Close() {
	syscall.Close(m.fd)
}
