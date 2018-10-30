package main

import (
	"fmt"
	"net"

	pb "github.com/asxcandrew/wbrkev/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = ":50051"
)

var chClients = make(chan interface{})

type Server struct{}

func (s *Server) GetCustomers(filter *pb.CustomerRequest, stream pb.Ingestor_GetCustomersServer) error {
	for c := range chClients {
		customer := c.(*pb.CustomerResponse)
		if err := stream.Send(customer); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) PushCustomer(context context.Context, in *pb.CustomerResponse) (*pb.StatusRequest, error) {
	chClients <- in
	return &pb.StatusRequest{Status: "OK"}, nil
}

func main() {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterIngestorServer(grpcServer, &Server{})
	reflection.Register(grpcServer)

	fmt.Println("GRPC server established")

	grpcServer.Serve(listen)
}
