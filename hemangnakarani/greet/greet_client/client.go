package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/greet/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello ! I'm Client.")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c)
	doServerStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Strting To do Unary RPC...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "TOPO",
			LastName:  "Nakarani",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not Catch Response : %v", err)
	}

	log.Printf("Response Unary: %v", res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Server Streaming In RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "TOPO",
			LastName:  "Nakarani",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not Catch the Stream : %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error While Reading the stream: %v", err)
		}

		log.Printf("Response Streaming: %v", msg.GetResult())
	}

}
