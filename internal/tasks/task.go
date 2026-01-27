package tasks

import (
	"context"
	"net"
	"netsnitch/internal/scan"
)

type Task interface {
	Protocol() scan.Protocol
	Execute(ctx context.Context, input []net.IP, results chan<- scan.Result)
}

var tasks = make(map[scan.Protocol]Task)

func Register(t Task) {

	tasks[t.Protocol()] = t
}
