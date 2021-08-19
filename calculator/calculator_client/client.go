package main

import (
	"context"
	"io"
	"log"

	"github.com/remiazv/golang-grpc-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)
	doServerStreaming(c)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient){
	request := &calculatorpb.DecompositionRequest{
		Number: 120,
	}

	resStream, err := c.Decomposition(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Decomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Response from Decomposition: %v", msg.GetResult())
	}
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	request := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 3,
	}

	response, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Calculating RPC: %v", err)
	}

	log.Printf("Response: %v", response.Result)
}
