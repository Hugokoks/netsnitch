package tcp_stealth

import (
	"encoding/binary"
	"net"
)

// ---- TCP PSEUDO-HEADER (for Checksum) ----
//
// [0-3]   Source IP Address
// [4-7]   Destination IP Address
// [8]     Zero Byte (Reserved)
// [9]     Protocol (6 for TCP)
// [10-11] TCP Length (20 bytes for our SYN)
// [12+]   The actual 20-byte TCP Header from above
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
	// Process data in 16-bit chunks (word by word).
	// We use a uint32 for the sum to capture any "overflow" (carries)
	// that exceed the 16-bit limit during addition.
	for i := 0; i+1 < len(data); i += 2 {
		bytesOfTwo := binary.BigEndian.Uint16(data[i : i+2])
		sum += uint32(bytesOfTwo)
	}
	// Handle the "Odd Byte" case.
	// If the data length is odd, treat the last byte as the most significant
	// byte of a 16-bit word (padded with a zero byte at the end).
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	// The "Carry-Around" Loop (1's Complement Arithmetic).
	// In TCP checksumming, any bits that overflow past bit 16 must be
	// added back to the bottom 16 bits.
	// (sum >> 16) extracts the carry, (sum & 0xFFFF) keeps the bottom 16 bits.

	//[0000000000000001,0000000000000000]
	//[0000000000000000,0000000000000001]
	for (sum >> 16) > 0 {

		// Take the lower 16 bits and sum them with the upper 16 bits
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	//Binary negation
	return ^uint16(sum)
}
