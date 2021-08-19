package main

import (
	"context"
	"log"
	"net"

	"grcp-udemy/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, request *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	firstNumber := request.GetCalculating().GetFirstNumber()
	secondNumber := request.GetCalculating().GetSecondNumber()

	response := &calculatorpb.CalculatorResponse{
		Result: firstNumber + secondNumber,
	}

	return response, nil
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