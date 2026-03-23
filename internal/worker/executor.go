package worker

import (
	"os/exec"
	"context"
	"github.com/jarvis1897/task_orchestrator/internal/common"
)

func ExecuteTask(ctx context.Context, task *common.Task) (string, error) {
	cmd := exec.CommandContext(ctx, task.Command, task.Args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}