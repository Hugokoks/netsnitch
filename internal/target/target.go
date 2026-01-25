package target

import (
	"net"
)


type Target struct {
	IP   net.IP
	Port int
}