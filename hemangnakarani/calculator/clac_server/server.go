package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/calculator/calcpb"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

type server struct {
	calcpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calcpb.SumRequest) (*calcpb.SumResponse, error) {

	log.Printf("Request:%v\n", req)

	res := &calcpb.SumResponse{
		SumResult: req.FirstNumber + req.SecondNumber,
	}

	return res, nil

}

func (*server) PrimeNumberDecomp(req *calcpb.PrimeNumberDecompRequest, stream calcpb.CalculatorService_PrimeNumberDecompServer) error {

	log.Printf("Server Streaming Request:%v\n", req)

	number := req.Number
	divisor := int64(2)

	for number > 1 {

		if number%divisor == 0 {
			stream.Send(&calcpb.PrimeNumberDecompResponse{
				PrimeFactor: divisor,
			})

			number = number / divisor
		} else {
			divisor++
			log.Printf("Divisor incremented to %v!\n", divisor)
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calcpb.CalculatorService_ComputeAverageServer) error {

	log.Printf("Client Streaming : Func Invoked\n")

	result := int64(0)
	cnt := 0

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Average: float64(result) / float64(cnt),
			})
		}

		if err != nil {
			log.Fatalf("Error While Reading Client Stream: %v", err)
		}

		number := req.Number
		result += number
		cnt++
	}

}

func (*server) FindMaximum(stream calcpb.CalculatorService_FindMaximumServer) error {

	fmt.Printf("FindMaximum Stream Startes,will Return Stream of maximums...\n")

	maxi := int64(0)

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error While Recving Stream.. %v\n", err)
			return err
		}

		number := req.GetNumber()

		if maxi < number {

			maxi = number

			sendErr := stream.Send(&calcpb.FindMaximumResponse{
				Maximum: maxi,
			})

			if sendErr != nil {
				log.Fatalf("error while sending stream: %v\n", sendErr)
				return sendErr
			}
		}
	}

}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {

	fmt.Printf("Request to get Square Root: %v\n", req)

	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a Negative Number:%v", number),
		)
	}

	return &calcpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil

}

func (*server) SquareWithDeadline(ctx context.Context, req *calcpb.SquareWithDeadlineRequest) (*calcpb.SquareWithDeadlineResponse, error) {

	fmt.Println("SquareWithDeadline Hitted !!!!")

	for i := 0; i < 3; i++ {

		fmt.Printf("Error if: %v\n", ctx.Err())

		if ctx.Err() == context.DeadlineExceeded {

			fmt.Println("Client Cancelled Request!")
			return nil, status.Errorf(codes.Canceled, "Client Cancelled Request")
		}

		time.Sleep(1 * time.Second)

	}
	number := req.GetNumber()

	return &calcpb.SquareWithDeadlineResponse{
		NumberSquare: number * number,
	}, nil
}

func main() {

	fmt.Println("Server init...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed To Listen : %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := false

	if tls {
		certFile := "D:/golib/src/gRPC-Course/hemangnakarani/ssl/server.crt"
		keyFile := "D:/golib/src/gRPC-Course/hemangnakarani/ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)

		if sslErr != nil {
			log.Fatalf("failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)

	calcpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}

	fmt.Println("Server started...")
}
