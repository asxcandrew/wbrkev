package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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
	fmt.Println("Starting ingestor service...")
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

	go func() {
		for c := range server.ch {
			raw := c.([]string)

			customer := &Customer{
				Name:         raw[1],
				Email:        raw[2],
				MobileNumber: raw[3],
			}
			customer.adjustValues()

			server.pushCustomer(customer)
		}
	}()
}

func main() {
	file, err := os.Open(os.Getenv("FILE_PATH"))

	if err != nil {
		log.Fatalf("Can`t open file: %v", err)
	}

	fileDriver := drivers.Driver(inputDriver)
	err = fileDriver.Parse(file, server.ch)

	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) pushCustomer(customer *Customer) {
	fmt.Println(customer)
	pbCustomer := &pb.CustomerResponse{
		Name:         customer.Name,
		Email:        customer.Email,
		MobileNumber: customer.MobileNumber,
	}
	response, err := s.grpcClient.PushCustomer(context.Background(), pbCustomer)

	if err != nil {
		fmt.Printf("Failed to send customer: %v\n", err)
		return
	}

	fmt.Printf("Customer sent, response: %v\n", response)
}

func (c *Customer) adjustValues() {
	re := regexp.MustCompile("[0-9]+")
	digits := re.FindAllString(c.MobileNumber, -1)
	digits = append([]string{"+44"}, digits...)

	c.MobileNumber = strings.Join(digits[:], "")
}
