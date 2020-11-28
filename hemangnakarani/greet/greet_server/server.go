package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Printf("Greet Func Invoked %v", req)

	firstname := req.GetGreeting().GetFirstName()
	result := "Hello " + firstname
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Printf("Greet Func Invoked %v", req)
	firstname := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {

		result := "Hello " + firstname + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Printf("LongGreet Func Invoked with client stream req\n")

	result := "Hello "

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error While Reading Client Stream: %v", err)
		}

		firstname := req.GetGreeting().GetFirstName()
		result += firstname + "! "
	}

}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {

	fmt.Printf("Greeet Everyone BiDi Invoked\n")

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error While Reading Client Stream: %v", err)
			return err
		}

		result := "Hello " + req.GetGreeting().GetFirstName() + " !!"

		sendError := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if sendError != nil {
			log.Fatalf("Error While Sending to stream Client: %v", err)
			return sendError
		}

	}

}

func main() {

	fmt.Println("Hello World !!!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed To Listen : %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
