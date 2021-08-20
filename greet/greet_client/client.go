package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/remiazv/golang-grpc-course/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a client!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDirectionalStreaming(c)
	doUnaryCallWithDeadline(c, 1 * time.Second) // should complete
	doUnaryCallWithDeadline(c, 5 * time.Second) // should timeout
}

func doUnary(c greetpb.GreetServiceClient){
	log.Println("Starting to do a Unary RPC...")

	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName: "Doe",
		},
	}
	response, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	
	log.Printf("Response: %v", response.Result)
}

func doUnaryCallWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration){
	log.Println("Starting to do a UnaryWithDeadline RPC...")

	request := &greetpb.GreetWithDeadLineRequest{
		FirstName: "John",
		LastName: "Doe",
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := c.GreetWithDeadLine(ctx, request)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit!, Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("Error while calling UnaryWithDeadline RPC: %v", err)
		}
		return
	}
	
	log.Printf("Response: %v", response.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient){
	fmt.Println("Starting to do a Server Streaming RPC...")
	
	request := &greetpb.GreetManyTimesRequest{
		FirstName: "John",
		LastName: "Doe",
	}

	resStream, err := c.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
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

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
	

}

func doClientStreaming(c greetpb.GreetServiceClient){
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			FirstName: "John",
			LastName: "Doe",
		},
		{
			FirstName: "Clark",
			LastName: "Kent",
		},
		{
			FirstName: "Silvio",
			LastName: "Santos",
		},
		{
			FirstName: "Mary",
			LastName: "help",
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Erro while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Erro while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v", res)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient){
	fmt.Println("Starting to do a BiDi RPC...")

	requests := []*greetpb.GreetEveryoneRequest{
		{
			FirstName: "John",
			LastName: "Doe",
		},
		{
			FirstName: "Clark",
			LastName: "Kent",
		},
		{
			FirstName: "Silvio",
			LastName: "Santos",
		},
		{
			FirstName: "Mary",
			LastName: "help",
		},
	}
	
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Erro while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v", req)
			sendError := stream.Send(req)
			if sendError != nil {
				log.Fatalf("Erro while sending message: %v", sendError)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		// function to receive a bunch of messages
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Erro while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v", response.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}