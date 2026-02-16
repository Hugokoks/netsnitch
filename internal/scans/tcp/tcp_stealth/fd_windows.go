//go:build windows

package tcp_stealth

import "syscall"

type socketFD = syscall.Handle
