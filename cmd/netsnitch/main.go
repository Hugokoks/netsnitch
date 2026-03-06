package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"netsnitch/internal/engine"
	"netsnitch/internal/fingerprint"
	"netsnitch/internal/input"
	_ "netsnitch/internal/scans/arp"
	_ "netsnitch/internal/scans/tcp"
	_ "netsnitch/internal/scans/udp"
	"netsnitch/internal/tasks"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults")
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	if len(os.Args) < 2 {
		fmt.Println("usage: netsnitch <query>")
		os.Exit(1)
	}

	// init fingerprint engine
	fpEngine := fingerprint.NewEngine()
	// load fingerprint data
	if err := fpEngine.LoadRules("data/rules.json"); err != nil {
		log.Fatalf("detect.json load error: %v", err)
	}
	if err := fpEngine.LoadProbes("data/probes.json"); err != nil {
		log.Fatalf("probes.json load error: %v", err)
	}
	// Parse input
	query, err := input.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("input error:", err)
		os.Exit(1)
	}

	//Stage → Tasks
	var allTasks []tasks.Task

	for _, cfg := range query.Configs {
		////load fingerprint engine into cfg
		cfg.Fingerprint = fpEngine

		////build task
		stageTasks, err := tasks.Build(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		allTasks = append(allTasks, stageTasks...)
	}
	// Engine
	if err := engine.Run(ctx, allTasks); err != nil {
		fmt.Println("engine error:", err)
		os.Exit(1)
	}
}
