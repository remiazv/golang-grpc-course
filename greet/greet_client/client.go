package main

import (
	"context"
	"fmt"
	"grcp-udemy/greet/greetpb"
	"log"

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
	log.Printf("Created client: %f", c)

	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient){
	log.Println("Starting to do a Unary RPC...")

	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Remi",
			LastName: "Azevedo",
		},
	}
	response, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	
	log.Printf("Response: %v", response.Result)
}