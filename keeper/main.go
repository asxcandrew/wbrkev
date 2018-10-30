package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/asxcandrew/wbrkev/protos"
	"google.golang.org/grpc"
)

const (
	serverAddr = "172.27.0.5"
	grpcPort   = ":50051"
)

func main() {
	fmt.Println("Starting keeper service...")
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

	stream, err := client.GetCustomers(context.Background(), &pb.CustomerRequest{Id: 1})
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		// Receiving the stream of data
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}

		if err != nil {
			log.Println(err)
		}
	}
}
