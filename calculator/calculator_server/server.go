package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/remiazv/golang-grpc-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, request *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	firstNumber := request.GetFirstNumber()
	secondNumber := request.GetSecondNumber()

	response := &calculatorpb.SumResponse{
		Result: firstNumber + secondNumber,
	}

	return response, nil
}

func (*server) Decomposition(request *calculatorpb.DecompositionRequest, stream calculatorpb.CalculatorService_DecompositionServer) error {
	n := request.GetNumber()
	var k int64 = 2

	for n > 1{

		if (n % k) == 0 {
			log.Println(k)
			n = n / k

			response := &calculatorpb.DecompositionResponse{
				Result: k,
			}
			stream.Send(response)
			time.Sleep(1000 * time.Millisecond)
		} else {
			k = k + 1
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	log.Printf("Server is ONLINE!!!!!!!!!!!!!!!")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}