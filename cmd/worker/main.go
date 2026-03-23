package main

import (
	"context"
	"github.com/jarvis1897/task_orchestrator/internal/worker"
)

func main() {
	w := worker.NewWorker("localhost:50050", "localhost:0")
	w.Start(context.Background())
}

