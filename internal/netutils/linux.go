//go:build linux

package netutils

import "syscall"

type SocketFD = int

const ProtoICMP = syscall.IPPROTO_ICMP
