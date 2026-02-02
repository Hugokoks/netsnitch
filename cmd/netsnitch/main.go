package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"netsnitch/internal/domain"
	"netsnitch/internal/engine"
	_ "netsnitch/internal/scans/arp_active"
	_ "netsnitch/internal/scans/tcp"
	"netsnitch/internal/tasks"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if len(os.Args) < 2 {
		fmt.Println("usage: netsnitch <cidr>")
		os.Exit(1)
	}

	cidr := os.Args[1]

	cfg := domain.Config{
		Timeout: 10 * time.Second,
		Type: domain.ARP_ACTIVE,
	}

	taks := tasks.Build(cfg,cidr)

	

	if err := engine.Run(ctx, taks, 1); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
