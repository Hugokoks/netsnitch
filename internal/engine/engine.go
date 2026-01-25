package engine

import (
	"context"
	"fmt"
	"netsnitch/internal/output"
	"netsnitch/internal/scan"
	"netsnitch/internal/scanner"

	"netsnitch/internal/target"
)

func Run(ctx context.Context, cidr string,cfg scan.Config) error {
	fmt.Println("[engine] starting scan for:", cidr)

	/////Parse Ip's from cidr 192.168.0.0/24
	ips, err := scan.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("parse CIDR failed: %w", err)
	}
	
	////create targets
	targets := target.BuildTargets(cfg, ips)
	
	////create Scanner
	sc,err := scanner.New(cfg)
	
	if err != nil {
		return err
	}
	
	////Create Scheduler
	scheduler := NewScheduler(ctx,cfg.Concurrency,sc)
	
	///write out scan results
	go output.ConsumeResults(ctx,scheduler.Results())	
	
	///start scanning
	scheduler.RunScan(targets)

	fmt.Println("[engine] scan finished, targets:", len(targets))
	return nil
}