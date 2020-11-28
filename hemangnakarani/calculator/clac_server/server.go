package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/calculator/calcpb"
	"io"
	"log"
	"net"

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

func main() {

	fmt.Println("Server init...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed To Listen : %v", err)
	}

	s := grpc.NewServer()

	calcpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}

	fmt.Println("Server started...")
}
