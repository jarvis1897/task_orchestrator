package worker

import (
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/jarvis1897/task_orchestrator/gen/proto"
	"github.com/jarvis1897/task_orchestrator/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Worker struct {
	NodeId  string
	Address string
	Client  pb.OrchestratorClient
}

func NewWorker(masterAddress string, workerAddress string) *Worker {
	nodeId := uuid.New().String()
	conn, err := grpc.NewClient(masterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := pb.NewOrchestratorClient(conn)
	_, err = client.Register(context.Background(), &pb.RegisterRequest{NodeId: nodeId, Address: workerAddress})
	if err != nil {
		panic(err)
	}
	return &Worker{
		NodeId:  nodeId,
		Address: workerAddress,
		Client:  client,
	}
}

func (w *Worker) Start(ctx context.Context) {
	go SendHeartbeat(ctx, w.NodeId, w.Client, 5)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			resp, err := w.Client.PollTask(ctx, &pb.PollRequest{NodeId: w.NodeId})
			if err != nil {
				continue
			}
			if !resp.HasTask {
				time.Sleep(2 * time.Second)
				continue
			} else {
				output, err := ExecuteTask(ctx, &common.Task{
					ID:      resp.TaskId,
					Command: resp.Command,
					Args:    resp.Args,
				})
				if err != nil {
					w.Client.ReportResult(ctx, &pb.TaskResult{
						NodeId:  w.NodeId,
						TaskId:  resp.TaskId,
						Success: false,
						Error:   err.Error(),
						Output:  output,
					})
				} else {
					w.Client.ReportResult(ctx, &pb.TaskResult{
						NodeId:  w.NodeId,
						TaskId:  resp.TaskId,
						Success: true,
						Output:  output,
						Error:   "",
					})
				}
			}
		}
	}
}
