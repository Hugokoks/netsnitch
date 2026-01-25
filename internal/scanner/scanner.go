package scanner

import (
	"context"
	"fmt"
	"netsnitch/internal/scan"
	"netsnitch/internal/target"
)

type Scanner interface {
	Scan(ctx context.Context, t target.Target) scan.Result

}

func New(cfg scan.Config) (Scanner, error) {
	switch cfg.Type {
	case scan.TCP:
		return &TCPScanner{Timeout: cfg.Timeout}, nil
	case scan.UDP:
		return &UDPScanner{Timeout: cfg.Timeout}, nil
	
	default:
		return nil, fmt.Errorf("unknown scan type")
	}
}
