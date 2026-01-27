package main

import (
	"context"
	"fmt"
	"netsnitch/internal/engine"
	"netsnitch/internal/scan"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	cfg := scan.Config{

		Type:        scan.ARP_ACTIVE,
		Timeout:     300 * time.Millisecond,
		Concurrency: 200,
	}

	if err := engine.Run(ctx, cidr, cfg); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

}
