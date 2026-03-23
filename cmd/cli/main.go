package main

import (
	"context"
	"log"
	"os"

	pb "github.com/jarvis1897/task_orchestrator/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2{
		log.Fatal("usage: ./cli <command> [args...]")
	}
	command := os.Args[1]
	args := os.Args[2:]

	conn, err := grpc.NewClient("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to master: %v", err)
	}
	defer conn.Close()
	client := pb.NewOrchestratorClient(conn)
	resp, err := client.SubmitTask(context.Background(), &pb.SubmitTaskRequest{Command: command, Args: args})
	if err != nil {
		log.Fatalf("failed to submit task: %v", err)
	}
	log.Printf("Task submitted successfully: %v", resp)
}