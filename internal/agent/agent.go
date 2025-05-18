// Package agent содержит реализацию агента для выполнения задач.
package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/denis-gr/GOCACL_DISTRIBUTED/internal/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Worker используется для обозначения ошибки, когда элемент не найден.
func Worker(delayMs int64, grpcAddress string) {
	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrchestratorServiceClient(conn)

	for {
		nextRun := time.Now().Add(time.Duration(delayMs) * time.Millisecond)
		task := getTask(client)
		if task != nil {
			result := performTask(task)
			err := sendResult(client, result)
			if err != nil {
				log.Println("Error sending result:", err)
			}
		}
		time.Sleep(time.Until(nextRun))
	}
}

func getTask(client pb.OrchestratorServiceClient) *pb.Task {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.GetTask(ctx, &pb.Empty{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.NotFound {
				return nil
			}
			log.Printf("gRPC error: code = %s, message = %s", st.Code(), st.Message())
		} else {
			log.Println("Error getting task:", err)
			panic(err)
		}
		return nil
	}

	return response.Task
}

func performTask(task *pb.Task) *pb.TaskResultRequest {
	wait := time.Now().Add(time.Duration(task.OperationTime) * time.Millisecond)

	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 != 0 {
			result = task.Arg1 / task.Arg2
		} else {
			result = 0
		}
	default:
		result = 0
	}

	time.Sleep(time.Until(wait))

	return &pb.TaskResultRequest{
		Id:     task.Id,
		Result: result,
	}
}

func sendResult(client pb.OrchestratorServiceClient, result *pb.TaskResultRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.SendResult(ctx, result)
	if err != nil {
		return fmt.Errorf("error sending result: %w", err)
	}

	return nil
}
