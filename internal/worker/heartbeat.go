package worker

import (
	"context"
	"log"
	"time"

	pb "github.com/jarvis1897/task_orchestrator/gen/proto"
)

func SendHeartBeat(ctx context.Context, nodeId string, client pb.OrchestratorClient, interval int) {
	ticker := time.NewTicker(time.Duration(interval))
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			_, err :=client.Heartbeat(ctx, &pb.HeartbeatRequest{NodeId: nodeId})
			if err != nil {
				log.Printf("failed to send heartbeats: %v", err)
			}
		case <- ctx.Done():
			return
		}
	}
}