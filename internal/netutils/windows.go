//go:build windows

package netutils

import "syscall"

type SocketFD = syscall.Handle

const ProtoICMP = 1
