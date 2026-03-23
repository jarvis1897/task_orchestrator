package main

import (
	"github.com/jarvis1897/task_orchestrator/internal/master"
)

func main() {
	registry := master.NewRegistry()
	scheduler := master.NewScheduler(registry)
	go registry.StartDeadNodeChecker()

	master.StartMasterServer(":50050", registry, scheduler)
}