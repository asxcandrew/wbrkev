package main

import (
	"context"
	"fmt"

	pb "github.com/asxcandrew/wbrkev/protos"
	"google.golang.org/grpc"
)

const (
	serverAddr  = "172.27.0.5"
	grpcPort    = ":50051"
	inputDriver = "CSV"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(serverAddr+grpcPort, opts...)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		return
	}
	defer conn.Close()
	client := pb.NewIngestorClient(conn)

	response, err := client.PushCustomer(context.Background(), &pb.CustomerResponse{Name: "someone"})

	if err != nil {
		fmt.Printf("Failed to send customer: %v\n", err)
		return
	}
	fmt.Printf("response: %v\n", response)
}
