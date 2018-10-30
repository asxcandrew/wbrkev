package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/asxcandrew/wbrkev/ingestor/drivers"
	pb "github.com/asxcandrew/wbrkev/protos"
	"google.golang.org/grpc"
)

const (
	serverAddr  = "172.27.0.5"
	grpcPort    = ":50051"
	inputDriver = "CSV"
)

type Customer struct {
	ID           uint64
	Name         string
	Email        string
	MobileNumber string
}

type Server struct {
	grpcClient pb.IngestorClient
	ch         chan interface{}
}

var server Server

func init() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(serverAddr+grpcPort, opts...)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		return
	}
	defer conn.Close()

	server.ch = make(chan interface{})
	server.grpcClient = pb.NewIngestorClient(conn)
}

func main() {
	file, err := os.Open(os.Getenv("FILE_PATH"))

	if err != nil {
		log.Fatalf("Can`t open file: %v", err)
	}

	fileDriver := drivers.Driver(inputDriver)
	err = fileDriver.Parse(file, server.ch)

	for c := range server.ch {
		customer := c.(*Customer)

		server.pushCustomer(customer)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) pushCustomer(customer *Customer) {
	response, err := s.grpcClient.PushCustomer(context.Background(), &pb.CustomerResponse{Name: "someone"})

	if err != nil {
		fmt.Printf("Failed to send customer: %v\n", err)
		return
	}

	fmt.Printf("Customer sent, response: %v\n", response)
}
