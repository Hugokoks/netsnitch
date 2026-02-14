package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"netsnitch/internal/engine"
	"netsnitch/internal/input"
	_ "netsnitch/internal/scans/arp"
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
		fmt.Println("usage: netsnitch <query>")
		os.Exit(1)
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

		////build task
		stageTasks := tasks.Build(cfg)
		allTasks = append(allTasks, stageTasks...)
	}
	// Engine
	workers := 200
	if err := engine.Run(ctx, allTasks, workers); err != nil {
		fmt.Println("engine error:", err)
		os.Exit(1)
	}
}
