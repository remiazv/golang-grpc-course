package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBidirectionalStreaming(c)
}

func doBidirectionalStreaming(c calculatorpb.CalculatorServiceClient){
	fmt.Println("Starting to do a FindMaximim BiDi Streaming RPC...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	go func(){
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func(){
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
				break
			}
			maximum := response.GetMaximum()
			fmt.Printf("Received a new maximum of ...: %v\n", maximum)
		}
		close(waitc)
	}()
	<-waitc
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient){
	requests := []*calculatorpb.ComputeAverageRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Erro while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Erro while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v", res)
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
