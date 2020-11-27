package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/calculator/calcpb"
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

	c := calcpb.NewCalculatorServiceClient(conn)
	doUnary(c)
	doServerStreaming(c)

}

func doUnary(c calcpb.CalculatorServiceClient) {

	fmt.Println("Strting To do Unary RPC...")

	req := &calcpb.SumRequest{
		FirstNumber:  27,
		SecondNumber: 32,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not Catch Response : %v", err)
	}

	log.Printf("Response Unary: %v", res.SumResult)

}

func doServerStreaming(c calcpb.CalculatorServiceClient) {

	fmt.Println("Starting Server Streaming RPC...")

	req := &calcpb.PrimeNumberDecompRequest{
		Number: 27 * 32 * 49,
	}

	resStream, err := c.PrimeNumberDecomp(context.Background(), req)

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

		log.Printf("Response Streaming: %v", msg.GetPrimeFactor())
	}

}
