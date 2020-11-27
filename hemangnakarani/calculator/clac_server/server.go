package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/calculator/calcpb"
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

	log.Printf("Streaming Request:%v\n", req)

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
