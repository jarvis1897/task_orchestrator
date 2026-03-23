package master

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	pb "github.com/jarvis1897/task_orchestrator/gen/proto"
	"github.com/jarvis1897/task_orchestrator/internal/common"
	"google.golang.org/grpc"
)

type MasterServer struct {
	pb.UnimplementedOrchestratorServer
	registry *Registry
	scheduler *Scheduler
}

func (s *MasterServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// create a common.node from the request data
	node := &common.Node{
		ID: req.NodeId,
		Address: req.Address,
		LastHeartbeat: time.Now(),
		Status: "active",
	}

	s.registry.Register(node)
	log.Printf("node registered: %s at %s", node.ID, node.Address)

	return &pb.RegisterResponse{
		Success: true,
		Message: "registered successfully",
	}, nil
}

func (s *MasterServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.registry.Heartbeat(req.NodeId)
	return &pb.HeartbeatResponse{
		Alive: true,
	}, nil
}

func (s *MasterServer) PollTask(ctx context.Context, req *pb.PollRequest) (*pb.TaskResponse, error) {
	task := s.scheduler.PollTask(req.NodeId)
	if task == nil {
		return &pb.TaskResponse{
			HasTask: false,
		}, nil
	} 
	return &pb.TaskResponse{
		HasTask: true,
		TaskId: task.ID,
		Command: task.Command,
		Args: task.Args,
	}, nil
}

func (s *MasterServer) ReportResult(ctx context.Context, req *pb.TaskResult) (*pb.ReportResponse, error) {
	log.Printf("task result received: taskId=%s success=%v output=%s error=%s", req.TaskId, req.Success, req.Output, req.Error)
	return &pb.ReportResponse{
		Acknowledged: true,
	}, nil
}

func (s *MasterServer) SubmitTask(ctx context.Context, req *pb.SubmitTaskRequest) (*pb.SubmitTaskResponse, error) {
	task := &common.Task{
		ID: uuid.NewString(),
		Command: req.Command,
		Args: req.Args,
		Status: "pending",
	}
	s.scheduler.AddTask(task)
	return  &pb.SubmitTaskResponse{
		Success: true,
		Message: "task submitted successfully",
		TaskId: task.ID,
	}, nil
}

func StartMasterServer(address string, registry *Registry, scheduler *Scheduler) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterOrchestratorServer(grpcServer, &MasterServer{
		registry: registry,
		scheduler: scheduler,
	})

	log.Printf("master listening on %s", address)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}



