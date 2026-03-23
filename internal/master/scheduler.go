package master

import (
	"github.com/jarvis1897/task_orchestrator/internal/common"
)

type Scheduler struct {
	tasks chan *common.Task
	registry*Registry
}

func NewScheduler(registry *Registry) *Scheduler{
	return &Scheduler{
		tasks: make(chan *common.Task, 100),
		registry: registry,
	}
}

func (s *Scheduler) AddTask(task *common.Task) {
	s.tasks <- task
}

func (s *Scheduler) PollTask(nodeId string) *common.Task {
	select {
	case task := <- s.tasks:
		return task
	default:
		return nil
	}
}

