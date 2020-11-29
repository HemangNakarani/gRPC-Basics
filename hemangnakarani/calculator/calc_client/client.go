package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/calculator/calcpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello ! I'm Client.")

	tls := true
	opts := grpc.WithInsecure()

	if tls {

		certFile := "D:/golib/src/gRPC-Course/hemangnakarani/ssl/ca.crt"

		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")

		if sslErr != nil {
			log.Fatalf("Error While Loading CA trust Certificate: %v\n", sslErr)
		}

		opts = grpc.WithTransportCredentials(creds)

	}

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer conn.Close()

	c := calcpb.NewCalculatorServiceClient(conn)

	doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	//doErrorUnary(c)
	//doUnaryWithDeadline(c, 5*time.Second)
	//doUnaryWithDeadline(c, 1*time.Second)

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

func doClientStreaming(c calcpb.CalculatorServiceClient) {

	requests := []*calcpb.ComputeAverageRequest{
		&calcpb.ComputeAverageRequest{
			Number: 1,
		},
		&calcpb.ComputeAverageRequest{
			Number: 2,
		},
		&calcpb.ComputeAverageRequest{
			Number: 3,
		},
		&calcpb.ComputeAverageRequest{
			Number: 4,
		},
		&calcpb.ComputeAverageRequest{
			Number: 5,
		},
		&calcpb.ComputeAverageRequest{
			Number: 6,
		},
		&calcpb.ComputeAverageRequest{
			Number: 7,
		},
		&calcpb.ComputeAverageRequest{
			Number: 8,
		},
		&calcpb.ComputeAverageRequest{
			Number: 9,
		},
		&calcpb.ComputeAverageRequest{
			Number: 10,
		},
	}

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error While Calling ComputeAverage: %v", err)
	}

	for _, req := range requests {

		fmt.Printf("Sending Req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error While reciving resp from ComputeAverage: %v", err)
	}

	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doBiDiStreaming(c calcpb.CalculatorServiceClient) {

	requests := []*calcpb.FindMaximumRequest{
		&calcpb.FindMaximumRequest{
			Number: 1,
		},
		&calcpb.FindMaximumRequest{
			Number: 2,
		},
		&calcpb.FindMaximumRequest{
			Number: 3,
		},
		&calcpb.FindMaximumRequest{
			Number: 2,
		},
		&calcpb.FindMaximumRequest{
			Number: 1,
		},
		&calcpb.FindMaximumRequest{
			Number: 6,
		},
		&calcpb.FindMaximumRequest{
			Number: 7,
		},
		&calcpb.FindMaximumRequest{
			Number: 8,
		},
		&calcpb.FindMaximumRequest{
			Number: 6,
		},
		&calcpb.FindMaximumRequest{
			Number: 100,
		},
		&calcpb.FindMaximumRequest{
			Number: 1,
		},
		&calcpb.FindMaximumRequest{
			Number: 2,
		},
		&calcpb.FindMaximumRequest{
			Number: 3,
		},
		&calcpb.FindMaximumRequest{
			Number: 2,
		},
		&calcpb.FindMaximumRequest{
			Number: 1,
		},
		&calcpb.FindMaximumRequest{
			Number: 6,
		},
		&calcpb.FindMaximumRequest{
			Number: 7,
		},
		&calcpb.FindMaximumRequest{
			Number: 8,
		},
		&calcpb.FindMaximumRequest{
			Number: 6,
		},
		&calcpb.FindMaximumRequest{
			Number: 1000,
		},
	}

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error While Calling FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	go func() {

		for _, req := range requests {

			err := stream.Send(req)

			if err != nil {
				log.Fatalf("error while sending stream: %v\n", err)
			}

			fmt.Printf("Send: %v\n", req)
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

			fmt.Printf("Recieved.................................. : %v\n", res.GetMaximum())
		}

		close(waitc)

	}()

	<-waitc

}

func doErrorUnary(c calcpb.CalculatorServiceClient) {

	fmt.Println("Strting To do Unary RPC...")

	res, err := c.SquareRoot(context.Background(), &calcpb.SquareRootRequest{Number: -36})

	if err != nil {
		respError, ok := status.FromError(err)

		if ok {
			// actual error from gRPC - User Error
			fmt.Println(respError.Message())
			fmt.Println(respError.Code())
		} else {
			log.Fatalf("Error From Calling Function: %v", err)
		}

		return
	}

	fmt.Printf("Response from Server: %v\n", res.GetNumberRoot())
}

func doUnaryWithDeadline(c calcpb.CalculatorServiceClient, timeout time.Duration) {

	fmt.Println("Strting To do Unary with Deadline RPC...")

	req := &calcpb.SquareWithDeadlineRequest{
		Number: 5,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.SquareWithDeadline(ctx, req)

	if err != nil {

		respError, ok := status.FromError(err)

		if ok {
			// actual error from gRPC - User Error
			fmt.Println(respError.Message())
			fmt.Println(respError.Code())

			if respError.Code() == codes.DeadlineExceeded {
				fmt.Println("We Cancled request")
			}
		} else {
			log.Fatalf("Error From Calling Function: %v", err)
		}
		return
	}

	log.Printf("Response Unary: %v", res.GetNumberSquare())

}
