package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/greet/greetpb"
	"io"
	"log"
	"time"

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
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)

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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Client Streaming In RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hemang",
				LastName:  "---",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Het",
				LastName:  "---",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Harshad",
				LastName:  "---",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Himanshu",
				LastName:  "---",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Harsh",
				LastName:  "---",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error While Calling LongGreet: %v", err)
	}

	for _, req := range requests {

		fmt.Printf("Sending Req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error While reciving resp from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	fmt.Printf("Starting BiDi Streaming In RPC...\n")

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hemang",
				LastName:  "---",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Het",
				LastName:  "---",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Harshad",
				LastName:  "---",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Himanshu",
				LastName:  "---",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Harsh",
				LastName:  "---",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while Reading stream\n")
		return
	}

	waitc := make(chan struct{})

	go func() {

		for _, req := range requests {

			fmt.Printf("Sending Message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()

	}()

	go func() {

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while Reading stream\n")
				break
			}

			fmt.Printf("Recieved : %v\n", res.GetResult())
		}

		close(waitc)

	}()

	<-waitc
}
