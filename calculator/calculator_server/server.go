package main

import (
	"context"
	"fmt"
	"io"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Compute Average is online!!!")
	var result int64 = 0
	var count float64 = 0
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			response := float64(result) / count
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: response,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client: %v", err)
		}

		result += request.GetNumber()
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum was invoked\n")
	var maximum int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client: %v", err)
		}

		number := req.GetNumber()
		if number > maximum {
			maximum = number
			sendError := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: maximum,
			})
			if sendError != nil {
				log.Fatalf("Error while sending data to client: %v", err)
				return err
			}
		}
	}
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