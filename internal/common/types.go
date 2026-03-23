package common

import "time"

const (
	StatusPending = "pending"
	StatusRunning = "running"
	StatusDone    = "done"
	StatusFailed  = "failed"
)

type Task struct {
	ID      string
	Command string
	Args    []string
	Status  string //(e.g. "pending", "running", "done", "failed")
}

type Node struct {
	ID            string
	Address       string
	LastHeartbeat time.Time
	Status        string //(e.g. "active", "dead")
}
