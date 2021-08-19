package main

import (
	"context"
	"grcp-udemy/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient){
	request := &calculatorpb.CalculatorRequest{
		Calculating: &calculatorpb.Calculating{
			FirstNumber: 10,
			SecondNumber: 3,
		},
	}

	response, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Calculating RPC: %v", err)
	}

	log.Printf("Response: %v", response.Result)
}