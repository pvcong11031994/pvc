package main

import (
	"context"
	"io"
	"log"

	pb "demo/customer"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

// createCustomer calls the RPC method CreateCustomer of CustomerServer
func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {

	resp, err := client.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Fatalf("Could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new Customer has been added with id: %d", resp.Id)
	}
}

// getCustomers calls the RPC method GetCustomers of CustomerServer
func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {

	// calling the streaming API
	stream, err := client.GetCustomers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
	}
}

func main() {

	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Creates a new CustomerClient
	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    1,
		Name:  "CongPV 1",
		Email: "vancong1@gmail.com",
		Phone: "123456789",
		Addresses: []*pb.CustomerRequest_Address{
			{
				Street:            "111C Nguyen Lam",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: false,
			},
			{
				Street:            "111B Nguyen Lam",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: true,
			},
		},
	}

	// Create a new customer
	createCustomer(client, customer)

	customer = &pb.CustomerRequest{
		Id:    2,
		Name:  "CongPV 2",
		Email: "vancong2@gmail.com",
		Phone: "1234567890",
		Addresses: []*pb.CustomerRequest_Address{
			{
				Street:            "302 To Hien Thanh",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: true,
			},
		},
	}

	// Create a new customer
	createCustomer(client, customer)

	// Filter with an empty Keyword
	filter := &pb.CustomerFilter{Keyword: "CongPV 2"}
	getCustomers(client, filter)
}
