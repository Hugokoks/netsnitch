package scanner

import (
	"context"
	"netsnitch/internal/scan"
	"netsnitch/internal/target"
	"time"
)

type UDPScanner struct {
	Timeout time.Duration
}

func (s *UDPScanner ) Scan (ctx context.Context,t target.Target) scan.Result{



	return  scan.Result{}
}