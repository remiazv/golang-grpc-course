package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/remiazv/golang-grpc-course/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)
	doServerStreaming(c)
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