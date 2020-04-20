package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "demo-grpc-go/service"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	names := []*pb.HelloRequest{
		{Name: "1"}, {Name: "2"}, {Name: "3"},
		{Name: "4"}, {Name: "5"}, {Name: "6"},
	}
	stream, err := c.SayHellos(ctx)
	if err != nil {
		log.Fatalf("could not greets: %v", err)
	}
	wait := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(wait)
				return
			}
			if err != nil {
				log.Fatalf("Failed to greet: %v", err)
			}
			log.Printf("Greeting: %s", in.GetMessage())
		}
	}()
	for _, name1 := range names {
		if err := stream.Send(name1); err != nil {
			log.Fatalf("Failed to send a name: %v", err)
		}
	}
	stream.CloseSend()
	<-wait
}
